//go:build darwin

package platform

func NewInputController() InputController {
	return NewMacOSInputController()
}
