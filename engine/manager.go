package engine

// Manager Managers exist to create and handle behaviour.
type Manager struct {
}

// Register this manager in the engine. This is called by the engine
// when the system is attached.
func (manager *Manager) Register() {
}

// RunConcurrent If this manager is supported to run concurrently, custom concurrency
// function should be defined here
func (manager *Manager) RunConcurrent() {
}

// Update Called every update loop.
// dt is the time elapsed since last called
func (manager *Manager) Update(dt float64) {
}

// Unregister Called when this manager is detached and destroyed by the
// engine
func (manager *Manager) Unregister() {
}

// PostUpdate Called at the end of each loop.
func (manager *Manager) PostUpdate() {
}
