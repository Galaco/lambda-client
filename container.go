package main

import (
	keyvalues "github.com/galaco/KeyValues"
	"github.com/galaco/lambda-client/internal/config"
	"github.com/galaco/lambda-core/filesystem"
	"github.com/galaco/tinygametools"
)

// Container
type Container struct {
	Config     config.Config
	GameInfo   keyvalues.KeyValue
	Filesystem filesystem.IFileSystem
	Window     tinygametools.Window
	Keyboard   tinygametools.Keyboard
	Mouse      tinygametools.Mouse
}
