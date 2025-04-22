package server

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type wsJsonInputState struct {
	Dx      int            `json:"dx"`
	Dy      int            `json:"dy"`
	Buttons ButtonStateMap `json:"buttons"`
	ScrollY int            `json:"scroll_y"`
}

func ParseJsonInputState(message []byte) (Command, error) {
	var jsonState wsJsonInputState

	if err := json.Unmarshal(message, &jsonState); err != nil {
		return Command{Type: UnknownCommand}, fmt.Errorf("JSON unmarshal error: %w", err)
	}

	cmd := Command{
		Type: InputStateCommand,
		Data: CommandData{
			Dx:      jsonState.Dx,
			Dy:      jsonState.Dy,
			Buttons: jsonState.Buttons,
			ScrollY: jsonState.ScrollY,
		},
	}

	if cmd.Data.Buttons == nil {

		cmd.Data.Buttons = make(ButtonStateMap)

	}

	return cmd, nil
}
