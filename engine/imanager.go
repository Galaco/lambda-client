package engine

// IManager Generic game manager.
// Different systems should implement these methods
type IManager interface {
	Register()
	RunConcurrent()
	Update(float64)
	Unregister()
	PostUpdate()
}
