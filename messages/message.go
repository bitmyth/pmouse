package messages

type EventKind string

const (
	MouseMoveRelative = "move"
	KeyUp             = "keyUp"
	KeyDown           = "keyDown"
	ClickLeft         = "clickLeft"
	ClickRight        = "clickRight"
)

type Event struct {
	Kind EventKind `json:"kind"`
	X    int       `json:"x"`
	Y    int       `json:"y"`
	Key  string    `json:"key"`
}
