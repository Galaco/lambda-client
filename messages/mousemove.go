package messages

import (
	"github.com/galaco/tinygametools"
)

const TypeMouseMove = tinygametools.EventName("MouseMove")

// MouseMove event object for when mouse is moved
type MouseMove struct {
	XY [2]float64
}

func (message *MouseMove) Type() tinygametools.EventName {
	return TypeMouseMove
}

func (message *MouseMove) Message() interface{} {
	return message.XY
}
