package server

import (
	"log"
	"net/http"
	"remox/platform"
	"time"

	"github.com/gorilla/websocket"
)

const (
	ParserTypeJSON     = "json"
	ParserTypeProtobuf = "protobuf"
)

type WebSocketServer struct {
	upgrader    websocket.Upgrader
	controller  platform.InputController
	parser      MessageParseFunc
	accelConfig AccelerationConfig
	parserType  string
}

type clientState struct {
	moveState         clientMoveState
	lastButtonState   ButtonStateMap
	accumulatedScroll float64
}

type clientMoveState struct {
	lastDx             int
	lastDy             int
	consistentSignX    float64
	consistentSignY    float64
	consistentDuration time.Duration
	lastMoveTime       time.Time
}

func NewWebSocketServer(controller platform.InputController, parser MessageParseFunc, accelConfig AccelerationConfig, parserType string) *WebSocketServer {
	if parser == nil {
		log.Fatal("WebSocketServer requires a non-nil message parser")
	}

	if parserType != ParserTypeJSON && parserType != ParserTypeProtobuf {
		log.Fatalf("WebSocketServer received invalid parserType: %s", parserType)
	}
	return &WebSocketServer{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		controller:  controller,
		parser:      parser,
		accelConfig: accelConfig,
		parserType:  parserType,
	}
}

func (s *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	defer func() {
		if err := s.controller.Close(); err != nil {
			log.Printf("Error during controller cleanup for client: %v", err)
		}
		log.Println("Platform controller cleaned up for client.")
	}()
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	log.Println("WebSocket client connected")

	var state = clientState{
		moveState: clientMoveState{
			lastMoveTime: time.Now(),
		},
		lastButtonState:   make(ButtonStateMap),
		accumulatedScroll: 0.0,
	}

	for {
		messageType, rawMessage, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Println("WebSocket client disconnected gracefully or abruptly.")
			} else {
				log.Println("WebSocket read error:", err)
			}
			break
		}

		currentTime := time.Now()

		expectedMessageType := websocket.TextMessage
		if s.parserType == ParserTypeProtobuf {
			expectedMessageType = websocket.BinaryMessage
		}
		if messageType == expectedMessageType {
			cmd, err := s.parser(rawMessage)
			if err != nil {
				log.Println("Message parsing failed:", err)
				continue
			}

			if cmd.Type == InputStateCommand {
				deltaTime := currentTime.Sub(state.moveState.lastMoveTime)
				handleMoveCommand(s.controller, cmd.Data.Dx, cmd.Data.Dy, s.accelConfig, &state.moveState, deltaTime)
				state.moveState.lastMoveTime = currentTime
				handleButtonState(s.controller, cmd.Data.Buttons, state.lastButtonState)
				handleScroll(s.controller, cmd.Data.ScrollY, &state, s.accelConfig)
			} else {
				log.Printf("Received unexpected command type: %s", cmd.Type)
			}
		} else {

			log.Printf("Received unexpected WebSocket message type: %d (expected: %d)", messageType, expectedMessageType)

			continue
		}
	}
	log.Println("Exiting message loop for client.")
}
