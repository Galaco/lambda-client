package behaviour

import (
	"github.com/galaco/lambda-client/engine"
	"github.com/galaco/lambda-client/messages"
	"github.com/galaco/tinygametools"
)

// Closeable Simple struct to control engine shutdown utilising the internal event manager
type Closeable struct {
	target *engine.Engine
}

// CallbackMouseMove function will shutdown the engine
func (closer Closeable) ReceiveMessage(message tinygametools.Event) {
	if message.Type() == messages.TypeKeyDown {
		if message.(*messages.KeyDown).Key == tinygametools.KeyEscape {
			// Will shutdown the engine at the end of the current loop
			closer.target.Close()
		}
	}
}

func NewCloseable(target *engine.Engine) *Closeable {
	return &Closeable{
		target: target,
	}
}
