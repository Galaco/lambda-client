package behaviour

import (
	"github.com/galaco/Lambda-Client/input/keyboard"
	"github.com/galaco/Lambda-Client/messages"
	"github.com/galaco/Lambda-Core/core"
	"github.com/galaco/Lambda-Core/core/event"
)

// Closeable Simple struct to control engine shutdown utilising the internal event manager
type Closeable struct {
	target *core.Engine
}

// CallbackMouseMove function will shutdown the engine
func (closer Closeable) ReceiveMessage(message event.IMessage) {
	if message.Type() == messages.TypeKeyDown {
		if message.(*messages.KeyDown).Key == keyboard.KeyEscape {
			// Will shutdown the engine at the end of the current loop
			closer.target.Close()
		}
	}
}

func NewCloseable(target *core.Engine) *Closeable {
	return &Closeable{
		target: target,
	}
}
