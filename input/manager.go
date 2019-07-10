package input

import (
	"github.com/galaco/lambda-client/engine"
	"github.com/galaco/lambda-client/event"
	"github.com/galaco/lambda-client/input/keyboard"
	"github.com/galaco/lambda-client/messages"
	"github.com/galaco/tinygametools"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

// manager handles user input from mouse and keyboard
// in a specific window
type Manager struct {
	engine.Manager
	MouseCoordinates mgl64.Vec2
	window           *tinygametools.Window
	lockMouse        bool

	mouse    *tinygametools.Mouse
	keyboard *tinygametools.Keyboard
}

// Register prepares this manager to listen for events from a window
func (manager *Manager) Register() {
	manager.keyboard.AddKeyCallback(manager.KeyCallback)
	manager.mouse.AddMousePosCallback(manager.MouseCallback)

	manager.keyboard.RegisterCallbacks(manager.window)
	manager.mouse.RegisterCallbacks(manager.window)

	_ = event.Dispatcher().Subscribe(messages.TypeKeyDown, GetKeyboard().ReceiveMessage, GetKeyboard())
	_ = event.Dispatcher().Subscribe(messages.TypeKeyReleased, GetKeyboard().ReceiveMessage, GetKeyboard())
	_ = event.Dispatcher().Subscribe(messages.TypeMouseMove, GetMouse().CallbackMouseMove, GetMouse())
}

// Update prepares data constructs that represent mouse & keyboard state with
// updated information on the current input state.
func (manager *Manager) Update(dt float64) {
	// Get window size
	x, y := manager.window.Handle().GetSize()
	if GetKeyboard().IsKeyDown(keyboard.KeyE) {
		manager.lockMouse = true
		manager.window.Handle().SetCursorPos(float64(x)/2, float64(y)/2)
		manager.window.Handle().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	} else {
		manager.lockMouse = false
		manager.window.Handle().SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	}

	GetMouse().Update()
	glfw.PollEvents()
}
func (manager *Manager) PostUpdate() {
	GetMouse().PostUpdate()
}

// Unregister
func (manager *Manager) Unregister() {

}

// KeyCallback called whenever a key event occurs
func (manager *Manager) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		_ = event.Dispatcher().Publish(&messages.KeyDown{Key: tinygametools.Key(key)})
	case glfw.Repeat:
		_ = event.Dispatcher().Publish(&messages.KeyHeld{Key: tinygametools.Key(key)})
	case glfw.Release:
		_ = event.Dispatcher().Publish(&messages.KeyReleased{Key: tinygametools.Key(key)})
	}
}

// MouseCallback called whenever a mouse event occurs
func (manager *Manager) MouseCallback(window *glfw.Window, xpos float64, ypos float64) {
	if manager.lockMouse {
		manager.MouseCoordinates[0], manager.MouseCoordinates[1] = window.GetCursorPos()
		w, h := window.GetSize()
		_ = event.Dispatcher().Publish(&messages.MouseMove{
			XY: [2]float64{
				float64(float64(w/2) - manager.MouseCoordinates[0]),
				float64(float64(h/2) - manager.MouseCoordinates[1]),
			},
		})
		window.SetCursorPos(float64(w/2), float64(h/2))
	}
}

// NewInputManager
func NewInputManager(win *tinygametools.Window, mouse *tinygametools.Mouse, keyboard *tinygametools.Keyboard) *Manager {
	manager := &Manager{
		window:   win,
		mouse:    mouse,
		keyboard: keyboard,
	}

	return manager
}
