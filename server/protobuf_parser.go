package server

import (
	"fmt"

	pb "remox/proto"

	"google.golang.org/protobuf/proto"
)

func ParseProtobufInputState(message []byte) (Command, error) {
	var pbState pb.InputState

	if err := proto.Unmarshal(message, &pbState); err != nil {
		return Command{Type: UnknownCommand}, fmt.Errorf("protobuf unmarshal error: %w", err)
	}

	cmd := Command{
		Type: InputStateCommand,
		Data: CommandData{

			Dx:      int(pbState.GetDx()),
			Dy:      int(pbState.GetDy()),
			Buttons: pbState.GetButtons(),
			ScrollY: int(pbState.GetScrollY()),
		},
	}

	if cmd.Data.Buttons == nil {

		cmd.Data.Buttons = make(ButtonStateMap)

	}

	return cmd, nil
}
