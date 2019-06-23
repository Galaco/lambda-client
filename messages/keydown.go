package messages

import (
	"github.com/galaco/Lambda-Client/input/keyboard"
	"github.com/galaco/Lambda-Core/core/event"
)

const TypeKeyDown = event.MessageType("KeyDown")

// KeyDown event object for keydown
type KeyDown struct {
	event.Message
	Key keyboard.Key
}

// Type returns message type
func (message *KeyDown) Type() event.MessageType {
	return TypeKeyDown
}
