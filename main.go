package main

import (
	"fmt"
	"github.com/galaco/KeyValues"
	"github.com/galaco/Lambda-Client/behaviour"
	"github.com/galaco/Lambda-Client/behaviour/controllers"
	"github.com/galaco/Lambda-Client/internal/config"
	"github.com/galaco/Lambda-Client/internal/debug"
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

func main() {
	runtime.LockOSThread()

	//if err := debug.StartProfiling("profile"); err == nil {
	//	defer debug.StopProfiling()
	//}

	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Println("Handled panic:", recovered)
		}
	}()

	logger.SetWriter(debug.NewStdOut())
	logger.EnablePretty()
	// Load GameInfo.txt
	// GameInfo.txt includes fundamental properties about the game
	// and its resources locations
	cfg, err := config.Load("./config.json")
	if err != nil {
		logger.Panic(err)
	}
	gameInfo, err := gameinfo.LoadConfig(cfg.GameDirectory)
	if err != nil {
		logger.Panic(err)
	}

	// Register GameInfo.txt referenced resource paths
	// Filesystem module needs to know about all the possible resource
	// locations it can search.
	fs := filesystem.CreateFilesystemFromGameInfoDefinitions(cfg.GameDirectory, gameInfo)

	// Explicitly define fallbacks for missing resources
	// Defaults are defined, but if HL2 assets are not readable, then
	// the default may not be readable
	resource.Manager().SetErrorModelName("models/props/de_dust/du_antenna_A.mdl")
	resource.Manager().SetErrorTextureName("materials/error.vtf")

	// General engine setup
	Application := core.NewEngine()
	Application.Initialise()

	// Game specific setup
	gameName := SetGame(&game.CounterstrikeSource{}, gameInfo)

	Application.AddManager(&window.Manager{
		Name: gameName,
	})
	Application.AddManager(&renderer.Manager{})
	Application.AddManager(&controllers.Camera{})

	RegisterShutdownMethod(Application)

	scene.LoadFromFile(cfg.GameDirectory + cfg.Map, fs)

	// Start
	Application.SetSimulationSpeed(10)
	Application.Run()

	defer resource.Manager().Empty()
}

// SetGame registers game entities and returns game name
func SetGame(proj game.IGame, gameInfo *keyvalues.KeyValue) string {
	windowName := "Lambda-Client: A BSP Viewer"
	gameInfoNode, _ := gameInfo.Find("GameInfo")
	if gameInfoNode == nil {
		logger.Panic("gameinfo was not found.")
		return windowName
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
