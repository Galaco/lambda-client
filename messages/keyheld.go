package messages

import (
	"github.com/galaco/tinygametools"
)

const TypeKeyHeld = tinygametools.EventName("KeyHeld")

// KeyHeld event object for when key is held down
type KeyHeld struct {
	Key tinygametools.Key
}

func (message *KeyHeld) Type() tinygametools.EventName {
	return TypeKeyHeld
}

func (message *KeyHeld) Message() interface{} {
	return message.Key
}
