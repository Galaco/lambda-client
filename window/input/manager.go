package input

import (
	"github.com/galaco/Lambda-Client/input"
	"github.com/galaco/Lambda-Client/input/keyboard"
	"github.com/galaco/Lambda-Client/messages"
	"github.com/galaco/Lambda-Core/core/event"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

// manager handles user input from mouse and keyboard
// in a specific window
type Manager struct {
	MouseCoordinates mgl64.Vec2
	window           *glfw.Window
	lockMouse        bool
}

// Register prepares this manager to listen for events from a window
func (manager *Manager) Register(window *glfw.Window) {
	manager.window = window
	window.SetKeyCallback(manager.KeyCallback)
	window.SetCursorPosCallback(manager.MouseCallback)

	event.Manager().Listen(messages.TypeKeyDown, input.GetKeyboard().ReceiveMessage)
	event.Manager().Listen(messages.TypeKeyReleased, input.GetKeyboard().ReceiveMessage)
	event.Manager().Listen(messages.TypeMouseMove, input.GetMouse().CallbackMouseMove)
}

// Update prepares data constructs that represent mouse & keyboard state with
// updated information on the current input state.
func (manager *Manager) Update(dt float64) {
	// Get window size
	x, y := manager.window.GetSize()
	if input.GetKeyboard().IsKeyDown(keyboard.KeyE) {
		manager.lockMouse = true
		manager.window.SetCursorPos(float64(x)/2, float64(y)/2)
		manager.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	} else {
		manager.lockMouse = false
		manager.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	}

	input.GetMouse().Update()
	glfw.PollEvents()
}
func (manager *Manager) PostUpdate() {
	input.GetMouse().PostUpdate()
}

// Unregister
func (manager *Manager) Unregister() {

}

// KeyCallback called whenever a key event occurs
func (manager *Manager) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		event.Manager().Dispatch(&messages.KeyDown{Key: keyboard.Key(key)})
	case glfw.Repeat:
		event.Manager().Dispatch(&messages.KeyHeld{Key: keyboard.Key(key)})
	case glfw.Release:
		event.Manager().Dispatch(&messages.KeyReleased{Key: keyboard.Key(key)})
	}
}

// MouseCallback called whenever a mouse event occurs
func (manager *Manager) MouseCallback(window *glfw.Window, xpos float64, ypos float64) {
	if manager.lockMouse {
		manager.MouseCoordinates[0], manager.MouseCoordinates[1] = window.GetCursorPos()
		w, h := window.GetSize()
		event.Manager().Dispatch(&messages.MouseMove{
			X: float64(float64(w/2) - manager.MouseCoordinates[0]),
			Y: float64(float64(h/2) - manager.MouseCoordinates[1]),
		})
		window.SetCursorPos(float64(w/2), float64(h/2))
	}
}
