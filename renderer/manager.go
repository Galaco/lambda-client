package renderer

import (
	"github.com/galaco/Lambda-Client/engine"
	"github.com/galaco/Lambda-Client/renderer/cache"
	"github.com/galaco/Lambda-Client/renderer/gl"
	"github.com/galaco/Lambda-Client/renderer/ui"
	"github.com/galaco/Lambda-Client/scene"
	"github.com/galaco/lambda-core/vgui"
	"github.com/galaco/tinygametools"
	"strings"
	"sync"
)

type Manager struct {
	engine.Manager

	window *tinygametools.Window

	renderer   IRenderer
	uiRenderer ui.Renderer

	dynamicPropCache cache.PropCache
}

var cacheMutex sync.Mutex

func (manager *Manager) Register() {
	manager.renderer = gl.NewRenderer()

	manager.renderer.LoadShaders()

	//manager.uiRenderer.InitRenderContext(manager.window)
}

func (manager *Manager) Update(dt float64) {
	currentScene := scene.Get()

	if manager.dynamicPropCache.NeedsRecache() {
		manager.RecacheEntities(currentScene)
	}

	currentScene.CurrentCamera().Update(dt)
	currentScene.GetWorld().TestVisibility(currentScene.CurrentCamera().Transform().Position)

	renderableWorld := currentScene.GetWorld()

	// Begin actual rendering
	manager.renderer.StartFrame(currentScene.CurrentCamera())

	// Start with sky
	manager.renderer.DrawSkybox(renderableWorld.Sky())

	// Draw static world first
	manager.renderer.DrawBsp(renderableWorld)

	// Dynamic objects
	cacheMutex.Lock()
	for _, entry := range *manager.dynamicPropCache.All() {
		manager.renderer.DrawModel(entry.Model, entry.Transform.TransformationMatrix())
	}
	cacheMutex.Unlock()

	//manager.uiRenderer.DrawUI()

	manager.renderer.EndFrame()
}

func (manager *Manager) SetUIMasterPanel(panel *vgui.MasterPanel) {
	manager.uiRenderer.SetMasterPanel(panel)
}

func (manager *Manager) RecacheEntities(scene *scene.Scene) {
	c := cache.NewPropCache()
	go func() {
		for _, ent := range *scene.GetAllEntities() {
			if ent.KeyValues().ValueForKey("model") == "" {
				continue
			}
			m := ent.KeyValues().ValueForKey("model")
			// Its a brush entity
			if !strings.HasSuffix(m, ".mdl") {
				continue
			}
			// Its a point entity
			c.Add(ent)
		}

		cacheMutex.Lock()
		manager.dynamicPropCache = *c
		cacheMutex.Unlock()
	}()
}

func NewRenderManager(win *tinygametools.Window) *Manager {
	return &Manager{
		window: win,
	}
}
