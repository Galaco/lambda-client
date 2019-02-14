package main

import (
	"github.com/galaco/Lambda-Client/behaviour"
	"github.com/galaco/Lambda-Client/behaviour/controllers"
	"github.com/galaco/Lambda-Client/config"
	"github.com/galaco/Lambda-Client/messages"
	"github.com/galaco/Lambda-Client/renderer"
	"github.com/galaco/Lambda-Client/scene"
	"github.com/galaco/Lambda-Client/window"
	"github.com/galaco/Lambda-Core/core"
	"github.com/galaco/Lambda-Core/core/event"
	"github.com/galaco/Lambda-Core/core/filesystem"
	"github.com/galaco/Lambda-Core/core/logger"
	"github.com/galaco/Lambda-Core/core/resource"
	"github.com/galaco/Lambda-Core/game"
	"github.com/galaco/Lambda-Core/lib/gameinfo"
	"runtime"
)

// Client
// Client provides a .bsp loading and rendering environment
// It provides full bsp loading, with props and materials.
// Visibility data is also used.
func main() {
	runtime.LockOSThread()

	logger.EnablePretty()
	// Load GameInfo.txt
	// GameInfo.txt includes fundamental properties about the game
	// and its resources locations
	cfg, err := config.Load("./config.json")
	if err != nil {
		logger.Fatal(err)
	}
	_, err = gameinfo.LoadConfig(cfg.GameDirectory)
	if err != nil {
		logger.Fatal(err)
	}

	// Register GameInfo.txt referenced resource paths
	// Filesystem module needs to know about all the possible resource
	// locations it can search.
	filesystem.RegisterGameResourcePaths(config.Get().GameDirectory, gameinfo.Get())

	// Explicity define fallbacks for missing resources
	// Defaults are defined, but if HL2 assets are not readable, then
	// the default may not be readable
	resource.Manager().SetErrorModelName("models/props/de_dust/du_antenna_A.mdl")
	resource.Manager().SetErrorTextureName("materials/error.vtf")

	// General engine setup
	Application := core.NewEngine()
	Application.Initialise()

	// Game specific setup
	gameName := SetGame(&game.CounterstrikeSource{})

	Application.AddManager(&window.Manager{
		Name: gameName,
	})
	Application.AddManager(&renderer.Manager{})
	Application.AddManager(&controllers.Camera{})

	// Register behaviour that needs to exist outside of game simulation & control
	RegisterShutdownMethod(Application)

	scene.LoadFromFile(config.Get().GameDirectory + "/maps/de_dust2.bsp")

	// Start
	Application.SetSimulationSpeed(10)
	Application.Run()

	defer resource.Manager().Empty()
}

// SetGame registers game entities and returns game name
func SetGame(proj game.IGame) string {
	windowName := "Gource"
	gameInfoNode, _ := gameinfo.Get().Find("GameInfo")
	if gameInfoNode == nil {
		logger.Fatal("gameinfo was not found.")
	}
	gameNode, _ := gameInfoNode.Find("game")
	if gameNode != nil {
		windowName, _ = gameNode.AsString()
	}

	proj.RegisterEntityClasses()

	return windowName
}

// RegisterShutdownMethod Implements a way of shutting down the engine
func RegisterShutdownMethod(app *core.Engine) {
	event.Manager().Listen(messages.TypeKeyDown, behaviour.NewCloseable(app).ReceiveMessage)
}
