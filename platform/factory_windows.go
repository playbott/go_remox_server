//go:build windows

package platform

func NewInputController() InputController {
	return NewWindowsInputController()
}
