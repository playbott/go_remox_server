package platform

import (
	"fmt"
	"log"
	"sync"
	"syscall"

	"github.com/gonutz/input"
	"github.com/gonutz/w32/v2"
)

var (
	user32Dll        = syscall.NewLazyDLL("user32.dll")
	procMouseEvent   = user32Dll.NewProc("mouse_event")
	mouseEventFWheel = 0x0800
	wheelDelta       = 120
)

var _ InputController = (*windowsInputController)(nil)

type windowsInputController struct {
	mu           sync.Mutex
	buttonStates map[string]bool
}

func NewWindowsInputController() InputController {
	return &windowsInputController{buttonStates: make(map[string]bool)}
}

func (wic *windowsInputController) MoveMouseBy(dx, dy int) error {
	err := input.MoveMouseBy(dx, dy)
	if err != nil {
		return err
	}
	return nil
}

func (wic *windowsInputController) SetButtonState(button string, pressed bool) error {
	wic.mu.Lock()
	defer wic.mu.Unlock()
	cur, exists := wic.buttonStates[button]
	x, y, ok := w32.GetCursorPos()
	if !ok {
		return fmt.Errorf("failed pos")
	}
	if !exists || cur != pressed {
		var err error
		switch button {
		case "left":
			if pressed {
				err = input.LeftDown(x, y)
			} else {
				err = input.LeftUp()
			}
		case "right":
			if pressed {
				err = input.RightDown(x, y)
			} else {
				err = input.RightUp()
			}
		case "middle":
			if pressed {
				err = input.MiddleDown(x, y)
			} else {
				err = input.MiddleUp()
			}
		default:
			return fmt.Errorf("unsupported:%s", button)
		}
		if err != nil {
			return err
		}
		wic.buttonStates[button] = pressed
	}
	return nil
}

func (wic *windowsInputController) Scroll(deltaY int) error {
	if deltaY == 0 {
		return nil
	}
	var scrollAmount int
	if deltaY < 0 {
		scrollAmount = wheelDelta
	} else {
		scrollAmount = -wheelDelta
	}
	_, _, lastErr := procMouseEvent.Call(uintptr(mouseEventFWheel), 0, 0, uintptr(scrollAmount), 0)
	if lastErr != nil && lastErr.Error() != "The operation completed successfully." {
		log.Printf("Error calling mouse_event for scroll: %v", lastErr)
		return fmt.Errorf("mouse_event call failed: %w", lastErr)
	}
	return nil
}

func (wic *windowsInputController) Close() error {
	log.Println("cleanup...")
	wic.mu.Lock()
	defer wic.mu.Unlock()
	var lastErr error
	for btn, pressed := range wic.buttonStates {
		if pressed {
			var err error
			switch btn {
			case "left":
				err = input.LeftUp()
			case "right":
				err = input.RightUp()
			case "middle":
				err = input.MiddleUp()
			}
			if err != nil {
				log.Printf("release err %s: %v", btn, err)
				lastErr = err
			}
			wic.buttonStates[btn] = false
		}
	}
	return lastErr
}
