package event

import "github.com/galaco/tinygametools"

type Name tinygametools.EventName

var manager *tinygametools.EventManager

// Dispatcher
func Dispatcher() *tinygametools.EventManager {
	if manager == nil {
		manager = tinygametools.NewEventManager()
	}
	return manager
}
