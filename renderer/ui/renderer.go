package ui

import (
	"github.com/galaco/lambda-core/vgui"
	"github.com/galaco/tinygametools"
	"github.com/inkyblackness/imgui-go"
)

type Renderer struct {
	masterPanel  *vgui.MasterPanel
	imguiImpl    *imguiGlfw3
	imguiContext *imgui.Context
}

func (renderer *Renderer) InitRenderContext(win *tinygametools.Window) {
	renderer.imguiContext = imgui.CreateContext(nil)
	renderer.imguiImpl = imguiGlfw3Init(win.Handle())
}

func (renderer *Renderer) SetMasterPanel(ui *vgui.MasterPanel) {
	renderer.masterPanel = ui
}

func (renderer *Renderer) DrawUI() {
	renderer.imguiImpl.NewFrame()
	imgui.Render()
	renderer.imguiImpl.Render(imgui.RenderedDrawData())
}
