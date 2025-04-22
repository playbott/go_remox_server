package platform

type InputController interface {
	MoveMouseBy(dx, dy int) error

	SetButtonState(button string, pressed bool) error

	Scroll(deltaY int) error

	Close() error
}
