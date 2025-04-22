package server

type CommandType string

const (
	InputStateCommand CommandType = "input_state"

	UnknownCommand CommandType = "unknown"
)

type ButtonStateMap map[string]bool

type CommandData struct {
	Dx      int
	Dy      int
	Buttons ButtonStateMap
	ScrollY int
}

type Command struct {
	Type CommandType
	Data CommandData
}

type MessageParseFunc func(message []byte) (Command, error)
