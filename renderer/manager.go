package renderer

import (
	"github.com/galaco/Lambda-Client/input"
	"github.com/galaco/Lambda-Client/input/keyboard"
	"github.com/galaco/Lambda-Client/renderer/cache"
	"github.com/galaco/Lambda-Client/renderer/gl"
	"github.com/galaco/Lambda-Client/scene"
	"github.com/galaco/Lambda-Core/core"
	"strings"
	"sync"
)

type Manager struct {
	core.Manager
	renderer IRenderer

	dynamicPropCache cache.PropCache
}

var cacheMutex sync.Mutex

func (manager *Manager) Register() {
	manager.renderer = gl.NewRenderer()

	manager.renderer.LoadShaders()
}

func (manager *Manager) Update(dt float64) {
	currentScene := scene.Get()

	if manager.dynamicPropCache.NeedsRecache() {
		manager.RecacheEntities(currentScene)
	}

	manager.updateRendererProperties()
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
		manager.renderer.DrawModel(entry.Model, entry.Transform.GetTransformationMatrix())
	}
	cacheMutex.Unlock()

	manager.renderer.EndFrame()
}

func (manager *Manager) updateRendererProperties() {
	manager.renderer.SetWireframeMode(input.GetKeyboard().IsKeyDown(keyboard.KeyX))
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