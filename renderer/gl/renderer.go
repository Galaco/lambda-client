package gl

import (
	"github.com/galaco/Lambda-Client/renderer/gl/bsp"
	material2 "github.com/galaco/Lambda-Client/renderer/gl/material"
	"github.com/galaco/Lambda-Client/renderer/gl/prop"
	"github.com/galaco/Lambda-Client/renderer/gl/shaders"
	"github.com/galaco/Lambda-Client/renderer/gl/shaders/sky"
	"github.com/galaco/Lambda-Client/scene/world"
	"github.com/galaco/Lambda-Core/core/entity"
	"github.com/galaco/Lambda-Core/core/event"
	"github.com/galaco/Lambda-Core/core/logger"
	"github.com/galaco/Lambda-Core/core/material"
	"github.com/galaco/Lambda-Core/core/mesh"
	"github.com/galaco/Lambda-Core/core/model"
	"github.com/galaco/Lambda-Core/core/resource/message"
	"github.com/galaco/gosigl"
	opengl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

//OpenGL renderer
type Renderer struct {
	lightmappedGenericShader gosigl.Context
	skyShader                gosigl.Context

	currentShaderId uint32

	uniformMap map[uint32]map[string]int32

	vertexDrawMode uint32

	matrixes struct {
		view       mgl32.Mat4
		projection mgl32.Mat4
	}
}

// Preparation function
// Loads shaders and sets necessary constants for opengls state machine
func (manager *Renderer) LoadShaders() {
	material2.TextureIdMap = map[string]gosigl.TextureBindingId{}
	prop.ModelIdMap = map[string][]*gosigl.VertexObject{}

	event.Manager().Listen(message.TypeTextureLoaded, material2.SyncTextureToGpu)
	event.Manager().Listen(message.TypeTextureUnloaded, material2.DestroyTextureOnGPU)
	event.Manager().Listen(message.TypeModelLoaded, prop.SyncPropToGpu)
	event.Manager().Listen(message.TypeModelUnloaded, prop.DestroyPropOnGPU)
	event.Manager().Listen(message.TypeMapLoaded, bsp.SyncMapToGpu)

	manager.lightmappedGenericShader = gosigl.NewShader()
	err := manager.lightmappedGenericShader.AddShader(shaders.Vertex, gosigl.VertexShader)
	if err != nil {
		logger.Fatal(err)
	}
	err = manager.lightmappedGenericShader.AddShader(shaders.Fragment, gosigl.FragmentShader)
	if err != nil {
		logger.Fatal(err)
	}
	manager.lightmappedGenericShader.Finalize()
	manager.skyShader = gosigl.NewShader()
	err = manager.skyShader.AddShader(sky.Vertex, opengl.VERTEX_SHADER)
	if err != nil {
		logger.Fatal(err)
	}
	err = manager.skyShader.AddShader(sky.Fragment, opengl.FRAGMENT_SHADER)
	if err != nil {
		logger.Fatal(err)
	}
	manager.skyShader.Finalize()

	//matrixes
	skyShaderMap := map[string]int32{}
	skyShaderMap["model"] = manager.skyShader.GetUniform("model")
	skyShaderMap["projection"] = manager.skyShader.GetUniform("projection")
	skyShaderMap["view"] = manager.skyShader.GetUniform("view")
	skyShaderMap["cubemapTexture"] = manager.lightmappedGenericShader.GetUniform("cubemapTexture")
	manager.uniformMap[manager.skyShader.Id()] = skyShaderMap

	manager.lightmappedGenericShader.UseProgram()
	lightmappedGenericShaderMap := map[string]int32{}
	lightmappedGenericShaderMap["model"] = manager.lightmappedGenericShader.GetUniform("model")
	lightmappedGenericShaderMap["projection"] = manager.lightmappedGenericShader.GetUniform("projection")
	lightmappedGenericShaderMap["view"] = manager.lightmappedGenericShader.GetUniform("view")
	//material properties
	lightmappedGenericShaderMap["albedoSampler"] = manager.lightmappedGenericShader.GetUniform("albedoSampler")
	lightmappedGenericShaderMap["useLightmap"] = manager.lightmappedGenericShader.GetUniform("useLightmap")
	lightmappedGenericShaderMap["lightmapTextureSampler"] = manager.lightmappedGenericShader.GetUniform("lightmapTextureSampler")
	manager.uniformMap[manager.lightmappedGenericShader.Id()] = lightmappedGenericShaderMap

	gosigl.SetLineWidth(32)
	gosigl.EnableBlend()
	gosigl.EnableDepthTest()
	gosigl.EnableCullFace(gosigl.Back, gosigl.WindingClockwise)

	gosigl.ClearColour(0, 0, 0, 1)
}

var numCalls = 0

// Called at the start of a frame
func (manager *Renderer) StartFrame(camera *entity.Camera) {
	manager.matrixes.projection = camera.ProjectionMatrix()
	manager.matrixes.view = camera.ViewMatrix()

	// Sky
	manager.skyShader.UseProgram()
	manager.setShader(manager.skyShader.Id())
	opengl.UniformMatrix4fv(manager.uniformMap[manager.skyShader.Id()]["projection"], 1, false, &manager.matrixes.projection[0])
	opengl.UniformMatrix4fv(manager.uniformMap[manager.skyShader.Id()]["view"], 1, false, &manager.matrixes.view[0])

	manager.lightmappedGenericShader.UseProgram()
	manager.setShader(manager.lightmappedGenericShader.Id())

	//matrixes
	opengl.UniformMatrix4fv(manager.uniformMap[manager.lightmappedGenericShader.Id()]["projection"], 1, false, &manager.matrixes.projection[0])
	opengl.UniformMatrix4fv(manager.uniformMap[manager.lightmappedGenericShader.Id()]["view"], 1, false, &manager.matrixes.view[0])

	gosigl.Clear(gosigl.MaskColourBufferBit, gosigl.MaskDepthBufferBit)
}

// Called at the end of a frame
func (manager *Renderer) EndFrame() {
	//if glError := opengl.GetError(); glError != opengl.NO_ERROR {
	//	logger.Error("error: %d\n", glError)
	//}
	//logger.Notice("Calls: %d", numCalls)
	numCalls = 0
}

// Draw the main bsp world
func (manager *Renderer) DrawBsp(world *world.World) {
	if bsp.MapGPUResource == nil {
		return
	}

	modelMatrix := mgl32.Ident4()
	opengl.UniformMatrix4fv(manager.uniformMap[manager.currentShaderId]["model"], 1, false, &modelMatrix[0])
	gosigl.BindMesh(bsp.MapGPUResource)
	//manager.BindMesh(world.Bsp().Mesh())
	for _, cluster := range world.VisibleClusters() {
		for _, face := range cluster.Faces {
			manager.DrawFace(&face)
		}
	}
	for _, face := range world.Bsp().DefaultCluster().Faces {
		manager.DrawFace(&face)
	}
	for _, cluster := range world.VisibleClusters() {
		for _, prop := range cluster.StaticProps {
			manager.DrawModel(prop.GetModel(), prop.Transform().GetTransformationMatrix())
		}
	}
	for _, prop := range world.Bsp().DefaultCluster().StaticProps {
		manager.DrawModel(prop.GetModel(), prop.Transform().GetTransformationMatrix())
	}
}

// Draw skybox (bsp model, staticprops, sky material)
func (manager *Renderer) DrawSkybox(sky *world.Sky) {
	if sky == nil || bsp.MapGPUResource == nil{
		return
	}

	if sky.GetVisibleBsp() != nil {
		modelMatrix := sky.Transform().GetTransformationMatrix()
		opengl.UniformMatrix4fv(manager.uniformMap[manager.currentShaderId]["model"], 1, false, &modelMatrix[0])

		gosigl.BindMesh(bsp.MapGPUResource)
		//manager.BindMesh(sky.GetVisibleBsp().Mesh())
		for _, cluster := range sky.GetClusterLeafs() {
			for _, face := range cluster.Faces {
				manager.DrawFace(&face)
			}
		}
		for _, cluster := range sky.GetClusterLeafs() {
			for _, prop := range cluster.StaticProps {
				manager.DrawModel(prop.GetModel(), prop.Transform().GetTransformationMatrix())
			}
		}
	}

	//manager.DrawSkyMaterial(sky.GetCubemap())
}

// Render a mesh and its submeshes/primitives
func (manager *Renderer) DrawModel(model *model.Model, transform mgl32.Mat4) {
	opengl.UniformMatrix4fv(manager.uniformMap[manager.currentShaderId]["model"], 1, false, &transform[0])
	modelBinding := prop.ModelIdMap[model.GetFilePath()]
	if modelBinding == nil {
		return
	}
	for idx, mesh := range model.GetMeshes() {
		// Missing materials will be flat coloured
		if mesh == nil || mesh.GetMaterial() == nil {
			// We need the fallback material
			continue
		}
		manager.BindMesh(mesh, modelBinding[idx])
		gosigl.DrawArray(0, len(mesh.Vertices())/3)

		numCalls++
	}
}

func (manager *Renderer) BindMesh(target mesh.IMesh, meshBinding *gosigl.VertexObject) {
	gosigl.BindMesh(meshBinding)
	//target.Bind()
	// $basetexture
	if target.GetMaterial() != nil {
		mat := target.GetMaterial().(*material.Material)
		opengl.Uniform1i(manager.uniformMap[manager.currentShaderId]["albedoSampler"], 0)
		gosigl.BindTexture2D(gosigl.TextureSlot(0), material2.TextureIdMap[mat.Textures.Albedo.GetFilePath()])

		if mat.BumpMapName != "" && mat.Textures.Normal != nil {
			opengl.Uniform1i(manager.uniformMap[manager.currentShaderId]["normalSampler"], 1)
			gosigl.BindTexture2D(gosigl.TextureSlot(1), material2.TextureIdMap[mat.Textures.Normal.GetFilePath()])
		}
	}
	// Bind lightmap texture if it exists
	if target.GetLightmap() != nil {
		opengl.Uniform1i(manager.uniformMap[manager.currentShaderId]["useLightmap"], 0) // lightmaps disabled
		opengl.Uniform1i(manager.uniformMap[manager.currentShaderId]["lightmapTextureSampler"], 1)
		//target.GetLightmap().Bind()
	} else {
		opengl.Uniform1i(manager.uniformMap[manager.currentShaderId]["useLightmap"], 0)
	}
}

func (manager *Renderer) DrawFace(target *mesh.Face) {
	// Skip materialless faces
	if target.Material() == nil {
		return
	}

	// $basetexture
	mat := target.Material().(*material.Material)
	opengl.Uniform1i(manager.uniformMap[manager.currentShaderId]["albedoSampler"], 0)
	gosigl.BindTexture2D(gosigl.TextureSlot(0), material2.TextureIdMap[mat.Textures.Albedo.GetFilePath()])

	if mat.BumpMapName != "" && mat.Textures.Normal != nil {
		opengl.Uniform1i(manager.uniformMap[manager.currentShaderId]["normalSampler"], 1)
		gosigl.BindTexture2D(gosigl.TextureSlot(1), material2.TextureIdMap[mat.Textures.Normal.GetFilePath()])
	}

	// Bind lightmap texture if it exists
	if target.IsLightmapped() == true {
		opengl.Uniform1i(manager.uniformMap[manager.currentShaderId]["useLightmap"], 0) // lightmaps disabled
		opengl.Uniform1i(manager.uniformMap[manager.currentShaderId]["lightmapTextureSampler"], 1)
		//target.Lightmap().Bind()
	} else {
		opengl.Uniform1i(manager.uniformMap[manager.currentShaderId]["useLightmap"], 0)
	}
	gosigl.DrawArray(int(target.Offset()), int(target.Length()))
}

// Render the sky material
func (manager *Renderer) DrawSkyMaterial(skybox *model.Model) {
	if skybox == nil {
		return
	}
	var oldCullFaceMode int32
	opengl.GetIntegerv(opengl.CULL_FACE_MODE, &oldCullFaceMode)
	var oldDepthFuncMode int32
	opengl.GetIntegerv(opengl.DEPTH_FUNC, &oldDepthFuncMode)

	opengl.CullFace(opengl.FRONT)
	opengl.DepthFunc(opengl.LEQUAL)
	opengl.DepthMask(false)

	manager.skyShader.UseProgram()
	manager.setShader(manager.skyShader.Id())
	opengl.UniformMatrix4fv(manager.uniformMap[manager.skyShader.Id()]["projection"], 1, false, &manager.matrixes.projection[0])
	opengl.UniformMatrix4fv(manager.uniformMap[manager.skyShader.Id()]["view"], 1, false, &manager.matrixes.view[0])

	//DRAW
	//skybox.GetMeshes()[0].Bind()
	//skybox.GetMeshes()[0].GetMaterial().Bind()
	opengl.Uniform1i(manager.uniformMap[manager.currentShaderId]["cubemapSampler"], 0)
	manager.DrawModel(skybox, mgl32.Ident4())

	// End
	opengl.DepthMask(true)
	opengl.CullFace(uint32(oldCullFaceMode))
	opengl.DepthFunc(uint32(oldDepthFuncMode))

	// Back to default shader
	manager.lightmappedGenericShader.UseProgram()
	manager.setShader(manager.lightmappedGenericShader.Id())
}

// Change the draw format.
func (manager *Renderer) SetWireframeMode(mode bool) {
	if mode == true {
		gosigl.SetVertexDrawMode(opengl.LINES)
	} else {
		gosigl.SetVertexDrawMode(gosigl.Triangles)
	}
}

func (manager *Renderer) setShader(shader uint32) {
	if manager.currentShaderId != shader {
		manager.currentShaderId = shader
	}
}

func (manager *Renderer) Unregister() {
	manager.skyShader.Destroy()
	manager.lightmappedGenericShader.Destroy()
}

func NewRenderer() *Renderer {
	r := Renderer{}
	r.SetWireframeMode(false)
	r.uniformMap = map[uint32]map[string]int32{}

	return &r
}