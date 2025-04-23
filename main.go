package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"remox/configs"
	"remox/platform"
	"remox/server"
	"time"
)

func serveTesterHTML(w http.ResponseWriter, r *http.Request) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current working directory: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fp := filepath.Join(cwd, viper.GetString("TEST_HTML_FILE"))
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		log.Printf("ERROR: HTML tester file not found at %s", fp)
		http.NotFound(w, r)
		return
	}
	log.Printf("Serving HTML tester from: %s", fp)
	http.ServeFile(w, r, fp)
}

func main() {
	configs.LoadMain()

	log.Println("Starting Mouse Control Server (Tunable AccumScroll)...")

	inputController := platform.NewWindowsInputController()
	log.Println("Using Windows Input Controller")

	var messageParser server.MessageParseFunc
	var selectedParserType string
	parserType := viper.GetString("MESSAGE_FORMAT")

	if parserType == "protobuf" {
		messageParser = server.ParseProtobufInputState
		selectedParserType = server.ParserTypeProtobuf
		log.Println("Using Protobuf Input State Parser")
	} else {
		// parserType = json
		messageParser = server.ParseJsonInputState
		selectedParserType = server.ParserTypeJSON
		log.Println("Using JSON Input State Parser (Default)")
	}

	accelConfig := configs.DefaultAccelerationConfig()

	log.Println("--- Acceleration & Scroll Config ---")
	log.Printf(" Velocity Sensitivity : %.4f", accelConfig.VelocitySensitivity)
	log.Printf(" Velocity Threshold   : %.2f", accelConfig.VelocityThreshold)
	log.Printf(" Max Velocity Factor  : %.2f", accelConfig.MaxVelocityFactor)
	log.Printf(" Duration Sensitivity : %.6f", accelConfig.DurationSensitivity)
	log.Printf(" Min Consistent Dur Ms: %d", accelConfig.MinConsistentDurationMs/time.Millisecond)
	log.Printf(" Max Tracked Duration Ms: %d", accelConfig.MaxTrackedDurationMs/time.Millisecond)
	log.Printf(" Max Duration Factor  : %.2f", accelConfig.MaxDurationFactor)
	log.Printf(" Reset Time Gap Ms    : %d", accelConfig.ResetTimeGapMs/time.Millisecond)
	log.Printf(" Scroll Sensitivity   : %.3f", accelConfig.ScrollSensitivity)
	log.Printf(" Max Total Factor     : %.2f", accelConfig.MaxTotalFactor)
	log.Println("----------------------------------")

	wsServer := server.NewWebSocketServer(inputController, messageParser, accelConfig, selectedParserType)
	fs := http.FileServer(http.Dir("."))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/" {

			fmt.Fprintf(w, "WebSocket Server.\nWS endpoint: /ws\nTester page: /test.html")

			return
		}

		fs.ServeHTTP(w, r)
	})
	http.HandleFunc("/test", serveTesterHTML)
	http.HandleFunc("/ws", wsServer.HandleConnections)

	host := viper.GetString("HOST")
	port := viper.GetString("PORT")
	log.Printf("Server listening on port %s", port)
	log.Printf("Access tester at http://%s:%s/test", host, port)
	log.Printf("WebSocket endpoint ws://%s:%s/ws", host, port)
	if err := http.ListenAndServe(host+":"+port, nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
