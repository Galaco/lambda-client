package gl

import (
	"github.com/galaco/gosigl"
	"github.com/galaco/lambda-client/core/entity"
	"github.com/galaco/lambda-client/core/event"
	"github.com/galaco/lambda-client/core/lib/util"
	"github.com/galaco/lambda-client/core/material"
	"github.com/galaco/lambda-client/core/mesh"
	"github.com/galaco/lambda-client/core/model"
	"github.com/galaco/lambda-client/core/resource/message"
	"github.com/galaco/lambda-client/renderer/camera"
	"github.com/galaco/lambda-client/renderer/gl/bsp"
	material2 "github.com/galaco/lambda-client/renderer/gl/material"
	"github.com/galaco/lambda-client/renderer/gl/prop"
	"github.com/galaco/lambda-client/renderer/gl/shaders"
	"github.com/galaco/lambda-client/renderer/gl/shaders/sky"
	"github.com/galaco/lambda-client/scene/world"
	opengl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

//OpenGL renderer
type Renderer struct {
	lightmappedGenericShader gosigl.Context
	skyShader                gosigl.Context

	currentShaderId uint32

	uniformMap map[uint32]map[string]int32

	materialCache *material2.Cache

	matrices struct {
		view       mgl32.Mat4
		projection mgl32.Mat4
	}

	activeCamera *entity.Camera
	viewFrustum  *camera.Frustum
}

// LoadShaders Loads shaders and sets necessary constants for opengls state machine
func (renderer *Renderer) LoadShaders() {
	renderer.materialCache = material2.NewCache()
	prop.ModelIdMap = map[string][]*gosigl.VertexObject{}

	event.Manager().Listen(message.TypeMaterialLoaded, renderer.materialCache.SyncTextureToGpu)
	event.Manager().Listen(message.TypeMaterialUnloaded, renderer.materialCache.DestroyTextureOnGPU)
	event.Manager().Listen(message.TypeModelLoaded, prop.SyncPropToGpu)
	event.Manager().Listen(message.TypeModelUnloaded, prop.DestroyPropOnGPU)
	event.Manager().Listen(message.TypeMapLoaded, bsp.SyncMapToGpu)

	renderer.lightmappedGenericShader = gosigl.NewShader()
	err := renderer.lightmappedGenericShader.AddShader(shaders.Vertex, gosigl.VertexShader)
	if err != nil {
		util.Logger().Panic(err)
	}
	err = renderer.lightmappedGenericShader.AddShader(shaders.Fragment, gosigl.FragmentShader)
	if err != nil {
		util.Logger().Panic(err)
	}
	renderer.lightmappedGenericShader.Finalize()
	renderer.skyShader = gosigl.NewShader()
	err = renderer.skyShader.AddShader(sky.Vertex, opengl.VERTEX_SHADER)
	if err != nil {
		util.Logger().Panic(err)
	}
	err = renderer.skyShader.AddShader(sky.Fragment, opengl.FRAGMENT_SHADER)
	if err != nil {
		util.Logger().Panic(err)
	}
	renderer.skyShader.Finalize()

	//matrices
	skyShaderMap := map[string]int32{}
	skyShaderMap["model"] = renderer.skyShader.GetUniform("model")
	skyShaderMap["projection"] = renderer.skyShader.GetUniform("projection")
	skyShaderMap["view"] = renderer.skyShader.GetUniform("view")
	skyShaderMap["cubemapTexture"] = renderer.lightmappedGenericShader.GetUniform("cubemapTexture")
	renderer.uniformMap[renderer.skyShader.Id()] = skyShaderMap

	renderer.lightmappedGenericShader.UseProgram()
	lightmappedGenericShaderMap := map[string]int32{}
	lightmappedGenericShaderMap["model"] = renderer.lightmappedGenericShader.GetUniform("model")
	lightmappedGenericShaderMap["projection"] = renderer.lightmappedGenericShader.GetUniform("projection")
	lightmappedGenericShaderMap["view"] = renderer.lightmappedGenericShader.GetUniform("view")
	//material properties
	lightmappedGenericShaderMap["albedoSampler"] = renderer.lightmappedGenericShader.GetUniform("albedoSampler")
	lightmappedGenericShaderMap["normalSampler"] = renderer.lightmappedGenericShader.GetUniform("normalSampler")
	lightmappedGenericShaderMap["useLightmap"] = renderer.lightmappedGenericShader.GetUniform("useLightmap")
	lightmappedGenericShaderMap["lightmapTextureSampler"] = renderer.lightmappedGenericShader.GetUniform("lightmapTextureSampler")
	renderer.uniformMap[renderer.lightmappedGenericShader.Id()] = lightmappedGenericShaderMap

	gosigl.SetLineWidth(32)
	gosigl.EnableBlend()
	gosigl.EnableDepthTest()
	gosigl.EnableCullFace(gosigl.Back, gosigl.WindingClockwise)

	gosigl.ClearColour(0, 0, 0, 1)
}

var numCalls = 0

// Called at the start of a frame
func (renderer *Renderer) StartFrame(cam *entity.Camera) {
	renderer.matrices.projection = cam.ProjectionMatrix()
	renderer.matrices.view = cam.ViewMatrix()
	renderer.activeCamera = cam
	renderer.viewFrustum = camera.FrustumFromCamera(renderer.activeCamera)

	// Sky
	renderer.skyShader.UseProgram()
	renderer.setShader(renderer.skyShader.Id())
	opengl.UniformMatrix4fv(renderer.uniformMap[renderer.skyShader.Id()]["projection"], 1, false, &renderer.matrices.projection[0])
	opengl.UniformMatrix4fv(renderer.uniformMap[renderer.skyShader.Id()]["view"], 1, false, &renderer.matrices.view[0])

	renderer.lightmappedGenericShader.UseProgram()
	renderer.setShader(renderer.lightmappedGenericShader.Id())

	//matrices
	opengl.UniformMatrix4fv(renderer.uniformMap[renderer.lightmappedGenericShader.Id()]["projection"], 1, false, &renderer.matrices.projection[0])
	opengl.UniformMatrix4fv(renderer.uniformMap[renderer.lightmappedGenericShader.Id()]["view"], 1, false, &renderer.matrices.view[0])

	gosigl.Clear(gosigl.MaskColourBufferBit, gosigl.MaskDepthBufferBit)
}

// Called at the end of a frame
func (renderer *Renderer) EndFrame() {
	//if glError := opengl.GetError(); glError != opengl.NO_ERROR {
	//	logger.Error("error: %d\n", glError)
	//}
	//logger.Notice("Calls: %d", numCalls)
	numCalls = 0
}

// Draw the main bsp world
func (renderer *Renderer) DrawBsp(world *world.World) {
	if bsp.MapGPUResource == nil {
		return
	}

	modelMatrix := mgl32.Ident4()
	opengl.UniformMatrix4fv(renderer.uniformMap[renderer.currentShaderId]["model"], 1, false, &modelMatrix[0])
	gosigl.BindMesh(bsp.MapGPUResource)

	renderClusters := make([]*model.ClusterLeaf, 0)

	for idx, cluster := range world.VisibleClusters() {
		// test cluster visibility for this frame
		if !renderer.viewFrustum.IsCuboidInFrustum(cluster.Mins, cluster.Maxs) {
			continue
		}
		renderClusters = append(renderClusters, world.VisibleClusters()[idx])
		for _, face := range cluster.Faces {
			renderer.DrawFace(&face)
		}
	}

	// Render objects that dont seem to belong to a cluster
	for _, face := range world.Bsp().DefaultCluster().Faces {
		renderer.DrawFace(&face)
	}
	for _, cluster := range renderClusters {
		// This is a performance cheat. We measure from the cluster origin for staticProp fades, rather than staticProp origin
		distToCluster := float32(math.Sqrt(
			math.Pow(float64(cluster.Origin.X()-renderer.activeCamera.Transform().Position.X()), 2) +
				math.Pow(float64(cluster.Origin.Y()-renderer.activeCamera.Transform().Position.Y()), 2) +
				math.Pow(float64(cluster.Origin.Z()-renderer.activeCamera.Transform().Position.Z()), 2)))

		for _, staticProp := range cluster.StaticProps {
			//  Skip render if staticProp is fully faded
			if staticProp.FadeMaxDistance() > 0 && distToCluster >= staticProp.FadeMaxDistance() {
				continue
			}
			renderer.DrawModel(staticProp.Model(), staticProp.Transform().TransformationMatrix())
		}
	}
	for _, prop := range world.Bsp().DefaultCluster().StaticProps {
		renderer.DrawModel(prop.Model(), prop.Transform().TransformationMatrix())
	}
}

// Draw skybox (bsp model, staticprops, sky material)
func (renderer *Renderer) DrawSkybox(sky *world.Sky) {
	if sky == nil || bsp.MapGPUResource == nil {
		return
	}

	if sky.GetVisibleBsp() != nil {
		modelMatrix := sky.Transform().TransformationMatrix()
		opengl.UniformMatrix4fv(renderer.uniformMap[renderer.currentShaderId]["model"], 1, false, &modelMatrix[0])

		gosigl.BindMesh(bsp.MapGPUResource)
		//renderer.BindMesh(sky.GetVisibleBsp().Mesh())
		for _, cluster := range sky.GetClusterLeafs() {
			for _, face := range cluster.Faces {
				renderer.DrawFace(&face)
			}
		}
		for _, cluster := range sky.GetClusterLeafs() {
			for _, prop := range cluster.StaticProps {
				renderer.DrawModel(prop.Model(), prop.Transform().TransformationMatrix())
			}
		}
	}

	//renderer.DrawSkyMaterial(sky.GetCubemap())
}

// Render a mesh and its submeshes/primitives
func (renderer *Renderer) DrawModel(model *model.Model, transform mgl32.Mat4) {
	opengl.UniformMatrix4fv(renderer.uniformMap[renderer.currentShaderId]["model"], 1, false, &transform[0])
	modelBinding := prop.ModelIdMap[model.FilePath()]
	if modelBinding == nil {
		return
	}
	for idx, mesh := range model.Meshes() {
		// Missing materials will be flat coloured
		if mesh == nil || mesh.Material() == nil {
			// We need the fallback material
			continue
		}
		renderer.BindMesh(mesh, modelBinding[idx])
		gosigl.DrawArray(0, len(mesh.Vertices())/3)

		numCalls++
	}
}

func (renderer *Renderer) BindMesh(target mesh.IMesh, meshBinding *gosigl.VertexObject) {
	gosigl.BindMesh(meshBinding)
	//target.Bind()
	// $basetexture
	if target.Material() != nil {
		mat := target.Material().(*material.Material)
		opengl.Uniform1i(renderer.uniformMap[renderer.currentShaderId]["albedoSampler"], 0)
		gosigl.BindTexture2D(gosigl.TextureSlot(0), renderer.materialCache.FetchCachedTexture(mat.Textures.Albedo.FilePath()))

		if mat.Textures.Normal != nil {
			opengl.Uniform1i(renderer.uniformMap[renderer.currentShaderId]["normalSampler"], 1)
			gosigl.BindTexture2D(gosigl.TextureSlot(1), renderer.materialCache.FetchCachedTexture(mat.Textures.Normal.FilePath()))
		}
	}
	// Bind lightmap texture if it exists
	//if target.GetLightmap() != nil {
	//	opengl.Uniform1i(renderer.uniformMap[renderer.currentShaderId]["useLightmap"], 0) // lightmaps disabled
	//	opengl.Uniform1i(renderer.uniformMap[renderer.currentShaderId]["lightmapTextureSampler"], 2)
	//	//target.GetLightmap().Bind()
	//} else {
	opengl.Uniform1i(renderer.uniformMap[renderer.currentShaderId]["useLightmap"], 0)
	//}
}

func (renderer *Renderer) DrawFace(target *mesh.Face) {
	// Skip materialless faces
	if target.Material() == nil {
		return
	}

	// $basetexture
	mat := target.Material().(*material.Material)
	opengl.Uniform1i(renderer.uniformMap[renderer.currentShaderId]["albedoSampler"], 0)
	gosigl.BindTexture2D(gosigl.TextureSlot(0), renderer.materialCache.FetchCachedTexture(mat.Textures.Albedo.FilePath()))

	if mat.Textures.Normal != nil {
		opengl.Uniform1i(renderer.uniformMap[renderer.currentShaderId]["normalSampler"], 1)
		gosigl.BindTexture2D(gosigl.TextureSlot(1), renderer.materialCache.FetchCachedTexture(mat.Textures.Normal.FilePath()))
	}

	// Bind lightmap texture if it exists
	//if target.IsLightmapped() {
	//	opengl.Uniform1i(renderer.uniformMap[renderer.currentShaderId]["useLightmap"], 0) // lightmaps disabled
	//	opengl.Uniform1i(renderer.uniformMap[renderer.currentShaderId]["lightmapTextureSampler"], 2)
	//	//target.Lightmap().Bind()
	//} else {
	opengl.Uniform1i(renderer.uniformMap[renderer.currentShaderId]["useLightmap"], 0)
	//}
	gosigl.DrawArray(int(target.Offset()), int(target.Length()))
}

// Render the sky material
func (renderer *Renderer) DrawSkyMaterial(skybox *model.Model) {
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

	renderer.skyShader.UseProgram()
	renderer.setShader(renderer.skyShader.Id())
	opengl.UniformMatrix4fv(renderer.uniformMap[renderer.skyShader.Id()]["projection"], 1, false, &renderer.matrices.projection[0])
	opengl.UniformMatrix4fv(renderer.uniformMap[renderer.skyShader.Id()]["view"], 1, false, &renderer.matrices.view[0])

	//DRAW
	//skybox.GetMeshes()[0].Bind()
	//skybox.GetMeshes()[0].GetMaterial().Bind()
	opengl.Uniform1i(renderer.uniformMap[renderer.currentShaderId]["cubemapSampler"], 0)
	renderer.DrawModel(skybox, mgl32.Ident4())

	// End
	opengl.DepthMask(true)
	opengl.CullFace(uint32(oldCullFaceMode))
	opengl.DepthFunc(uint32(oldDepthFuncMode))

	// Back to default shader
	renderer.lightmappedGenericShader.UseProgram()
	renderer.setShader(renderer.lightmappedGenericShader.Id())
}

func (renderer *Renderer) setShader(shader uint32) {
	if renderer.currentShaderId != shader {
		renderer.currentShaderId = shader
	}
}

func (renderer *Renderer) Unregister() {
	renderer.skyShader.Destroy()
	renderer.lightmappedGenericShader.Destroy()
}

func NewRenderer() *Renderer {
	return &Renderer{
		uniformMap: map[uint32]map[string]int32{},
	}
}
