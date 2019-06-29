package messages

import (
	"github.com/galaco/tinygametools"
)

const TypeKeyReleased = tinygametools.EventName("KeyReleased")

// KeyReleased event object for key released
type KeyReleased struct {
	Key tinygametools.Key
}

func (message *KeyReleased) Type() tinygametools.EventName {
	return TypeKeyReleased
}

func (message *KeyReleased) Message() interface{} {
	return message.Key
}
