package window

import (
	"github.com/galaco/lambda-client/engine"
	"github.com/galaco/tinygametools"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Manager is responsible for managing this games window. Understand
// that there is a distinction between the window and the renderer.
// This manager provides a window that a rendering context can be
// obtained from, and device input handling.
type Manager struct {
	engine.Manager
	window *tinygametools.Window

	Name string
}

// Register will create a new Window
func (manager *Manager) Register() {
}

// Update simply calls the input manager that uses this window
func (manager *Manager) Update(dt float64) {
}

// Unregister will end input listening and kill any window
func (manager *Manager) Unregister() {
	glfw.Terminate()
}

// PostUpdate is called at the end of an update loop.
// In this case it simply SwapBuffers the window, (to display updated window
// contents)
func (manager *Manager) PostUpdate() {
	manager.window.Handle().SwapBuffers()
}

// NewWindowManager
func NewWindowManager(win *tinygametools.Window) *Manager {
	return &Manager{
		window: win,
	}
}
