package messages

import (
	"github.com/galaco/Lambda-Core/core/event"
)

const TypeMouseMove = event.MessageType("MouseMove")

// MouseMove event object for when mouse is moved
type MouseMove struct {
	event.Message
	X float64
	Y float64
}

func (message *MouseMove) Type() event.MessageType {
	return TypeMouseMove
}
