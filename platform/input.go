package platform

type InputController interface {
	MoveMouseBy(dx, dy int) error

	SetButtonState(button string, pressed bool) error

	ScrollBy(dx, dy int) error

	Close() error
}
