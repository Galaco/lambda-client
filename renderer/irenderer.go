package renderer

import (
	"github.com/galaco/lambda-client/core/entity"
	"github.com/galaco/lambda-client/core/model"
	"github.com/galaco/lambda-client/scene/world"
	"github.com/go-gl/mathgl/mgl32"
)

type IRenderer interface {
	StartFrame(*entity.Camera)
	LoadShaders()
	DrawBsp(*world.World)
	DrawSkybox(*world.Sky)
	DrawModel(*model.Model, mgl32.Mat4)
	DrawSkyMaterial(*model.Model)
	EndFrame()
	Unregister()
}
