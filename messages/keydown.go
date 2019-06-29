package messages

import (
	"github.com/galaco/tinygametools"
)

const TypeKeyDown = tinygametools.EventName("KeyDown")

// KeyDown event object for keydown
type KeyDown struct {
	Key tinygametools.Key
}

// Type returns message type
func (message *KeyDown) Type() tinygametools.EventName {
	return TypeKeyDown
}

// Message
func (message *KeyDown) Message() interface{} {
	return message.Key
}
