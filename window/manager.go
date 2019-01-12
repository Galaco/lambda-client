package window

import (
	"github.com/galaco/Lambda-Client/config"
	"github.com/galaco/Lambda-Client/window/input"
	"github.com/galaco/Lambda-Client/window/window"
	"github.com/galaco/Lambda-Core/core"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Manager is responsible for managing this games window. Understand
// that there is a distinction between the window and the renderer.
// This manager provides a window that a rendering context can be
// obtained from, and device input handling.
type Manager struct {
	core.Manager
	window *glfw.Window
	input  input.Manager

	Name string
}

// Register will create a new Window
func (manager *Manager) Register() {
	manager.window = window.Create(config.Get().Video.Width, config.Get().Video.Height, manager.Name)
	manager.input.Register(manager.window)
}

// Update simply calls the input manager that uses this window
func (manager *Manager) Update(dt float64) {
	manager.input.Update(0)
}

// Unregister will end input listening and kill any window
func (manager *Manager) Unregister() {
	manager.input.Unregister()
	glfw.Terminate()
}

// PostUpdate is called at the end of an update loop.
// In this case it simply SwapBuffers the window, (to display updated window
// contents)
func (manager *Manager) PostUpdate() {
	manager.input.PostUpdate()
	manager.window.SwapBuffers()
}
