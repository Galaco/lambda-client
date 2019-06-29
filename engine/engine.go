package engine

import (
	"time"
)

// Engine Game engine
// Only 1 should be initialised
type Engine struct {
	Managers        []IManager
	running         bool
	simulationSpeed float64
}

// Initialise the engine, and attached managers
func (engine *Engine) Initialise() {
}

// Run the engine
func (engine *Engine) Run() {
	engine.running = true

	// Begin the event manager thread in the background
	//event.Manager().ProcessQueue()
	// Anything that runs concurrently can start now
	for _, manager := range engine.Managers {
		manager.RunConcurrent()
	}

	dt := 0.0
	startingTime := time.Now().UTC()

	for engine.running {
		for _, manager := range engine.Managers {
			manager.Update(dt)
		}

		for _, manager := range engine.Managers {
			manager.PostUpdate()
		}

		dt = (float64(time.Now().UTC().Sub(startingTime).Nanoseconds()/1000000) / 1000) * engine.simulationSpeed
		startingTime = time.Now().UTC()
	}

	for _, manager := range engine.Managers {
		manager.Unregister()
	}
}

// AddManager Adds a new manager to the engine
func (engine *Engine) AddManager(manager IManager) {
	engine.Managers = append(engine.Managers, manager)
	manager.Register()
}

// Close marks the engine to exit at the end of the current loop
func (engine *Engine) Close() {
	engine.running = false
}

// SetSimulationSpeed allows for speeding up and slowing down the game clock
func (engine *Engine) SetSimulationSpeed(multiplier float64) {
	if multiplier < 0 {
		return
	}
	engine.simulationSpeed = multiplier
}

// NewEngine returns a new engine instance
func NewEngine() *Engine {
	return &Engine{
		simulationSpeed: 1.0,
	}
}
