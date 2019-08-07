package ui

import (
	"github.com/galaco/lambda-client/core/filesystem"
	vguiCore "github.com/galaco/lambda-client/core/loader/vgui"
	"github.com/galaco/lambda-client/core/vgui"
	"github.com/galaco/lambda-client/engine"
	"github.com/galaco/tinygametools"
)

type Gui struct {
	engine.Manager
	window      *tinygametools.Window
	masterPanel vgui.MasterPanel
}

func (ui *Gui) Register() {

}

func (ui *Gui) Update(dt float64) {
	ui.Render()
}

func (ui *Gui) Render() {
	ui.masterPanel.Draw()
}

// LoadVGUIResource
func (ui *Gui) LoadVGUIResource(fs filesystem.IFileSystem, filename string) error {
	p, err := vguiCore.LoadVGUI(fs, filename)
	if err != nil {
		return err
	}
	ui.masterPanel.AddChild(p)

	return nil
}

func (ui *Gui) MasterPanel() *vgui.MasterPanel {
	return &ui.masterPanel
}

func NewGUIManager(win *tinygametools.Window) *Gui {
	return &Gui{
		window: win,
	}
}
