package main

import (
	keyvalues "github.com/galaco/KeyValues"
	"github.com/galaco/lambda-client/internal/config"
	"github.com/galaco/tinygametools"
	"github.com/golang-source-engine/filesystem"
)

// Container
type Container struct {
	Config     config.Config
	GameInfo   keyvalues.KeyValue
	Filesystem *filesystem.FileSystem
	Window     tinygametools.Window
	Keyboard   tinygametools.Keyboard
	Mouse      tinygametools.Mouse
}
