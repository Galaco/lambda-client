package main

import (
	"fmt"
	"github.com/galaco/KeyValues"
	"github.com/galaco/lambda-client/behaviour"
	"github.com/galaco/lambda-client/behaviour/controllers"
	"github.com/galaco/lambda-client/engine"
	"github.com/galaco/lambda-client/event"
	"github.com/galaco/lambda-client/input"
	"github.com/galaco/lambda-client/internal/config"
	"github.com/galaco/lambda-client/internal/debug"
	"github.com/galaco/lambda-client/messages"
	"github.com/galaco/lambda-client/renderer"
	"github.com/galaco/lambda-client/scene"
	"github.com/galaco/lambda-client/ui/dialogs"
	"github.com/galaco/lambda-client/window"
	"github.com/galaco/lambda-core/game"
	"github.com/galaco/lambda-core/lib/gameinfo"
	"github.com/galaco/lambda-core/lib/util"
	"github.com/galaco/lambda-core/resource"
	"github.com/galaco/tinygametools"
	"github.com/golang-source-engine/filesystem"
	"log"
	"runtime"
)

func main() {
	runtime.LockOSThread()

	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Println("Caught panic:", recovered)
			dialogs.ErrorMessage(fmt.Errorf("%s", recovered))
		}
	}()

	util.Logger().SetWriter(debug.NewStdOut())
	util.Logger().EnablePretty()

	container := Container{}

	// Load GameInfo.txt
	// GameInfo.txt includes fundamental properties about the game
	// and its resources locations
	cfg, err := config.Load("./config.json")
	if err != nil {
		util.Logger().Panic(err)
	}
	container.Config = *cfg
	// Load GameInfo
	gameInfo, err := gameinfo.LoadConfig(cfg.GameDirectory)
	if err != nil {
		util.Logger().Panic(err)
	}
	container.GameInfo = *gameInfo

	// Register GameInfo.txt referenced resource paths
	// Filesystem module needs to know about all the possible resource
	// locations it can search.
	container.Filesystem, err = filesystem.CreateFilesystemFromGameInfoDefinitions(cfg.GameDirectory, gameInfo, true)
	if err != nil {
		util.Logger().Error(err)
	} else {
		for _, loc := range container.Filesystem.EnumerateResourcePaths() {
			util.Logger().Notice(fmt.Sprintf("Loaded resource path: %s", loc))
		}
	}

	// Explicitly define fallbacks for missing resources
	// Defaults are defined, but if HL2 assets are not readable, then
	// the default may not be readable
	resource.Manager().SetErrorModelName("models/error.mdl")
	resource.Manager().SetErrorTextureName("materials/error.vtf")
	defer resource.Manager().Empty()

	// General engine setup
	Application := engine.NewEngine()
	Application.Initialise()

	// Game specific setup
	gameName := SetGame(&game.CounterstrikeSource{}, gameInfo)

	// Create window
	win, err := tinygametools.NewWindow(config.Get().Video.Width, config.Get().Video.Height, gameName)
	if err != nil {
		util.Logger().Panic(err)
	}
	win.Handle().MakeContextCurrent()
	container.Window = *win

	container.Mouse = *tinygametools.NewMouse()
	container.Keyboard = *tinygametools.NewKeyboard()

	Application.AddManager(window.NewWindowManager(&container.Window))
	Application.AddManager(input.NewInputManager(&container.Window, &container.Mouse, &container.Keyboard))
	//vgui := ui.NewGUIManager(&container.Window)
	renderManager := renderer.NewRenderManager(&container.Window)
	//renderManager.SetUIMasterPanel(vgui.MasterPanel())
	Application.AddManager(renderManager)
	Application.AddManager(&controllers.Camera{})

	//_ = vgui.LoadVGUIResource(container.Filesystem, "gamemenu")
	//vgui.MasterPanel().Resize(float64(cfg.Video.Width), float64(cfg.Video.Width))
	//Application.AddManager(vgui)

	RegisterShutdownMethod(Application)

	sceneName, loadSceneError := dialogs.OpenFile("Valve BSP files", "bsp")
	if loadSceneError != nil {
		dialogs.ErrorMessage(loadSceneError)
		log.Fatal()
	} else {
		scene.LoadFromFile(sceneName, container.Filesystem)
	}

	// Start
	Application.SetSimulationSpeed(10)
	Application.Run()
}

// SetGame registers game entities and returns game name
func SetGame(proj game.IGame, gameInfo *keyvalues.KeyValue) string {
	windowName := "lambda-client: A BSP Viewer"
	gameInfoNode, _ := gameInfo.Find("GameInfo")
	if gameInfoNode == nil {
		util.Logger().Panic("gameinfo was not found.")
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
func RegisterShutdownMethod(app *engine.Engine) {
	_ = event.Dispatcher().Subscribe(messages.TypeKeyDown, behaviour.NewCloseable(app).ReceiveMessage, app)
}
