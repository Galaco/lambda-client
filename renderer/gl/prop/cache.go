package prop

import (
	"github.com/galaco/gosigl"
	"github.com/galaco/lambda-client/core/event"
	"github.com/galaco/lambda-client/core/resource/message"
)

var ModelIdMap map[string][]*gosigl.VertexObject

func SyncPropToGpu(dispatched event.IMessage) {
	msg := dispatched.(*message.PropLoaded)
	if ModelIdMap[msg.Resource.FilePath()] != nil {
		return
	}
	vals := make([]*gosigl.VertexObject, len(msg.Resource.Meshes()))

	for idx, mesh := range msg.Resource.Meshes() {
		gpuObject := gosigl.NewMesh(mesh.Vertices())
		gosigl.CreateVertexAttribute(gpuObject, mesh.UVs(), 2)
		gosigl.CreateVertexAttribute(gpuObject, mesh.Normals(), 3)

		if len(mesh.Tangents()) == 0 {
			mesh.GenerateTangents()
		}
		gosigl.CreateVertexAttribute(gpuObject, mesh.Tangents(), 4)
		gosigl.CreateVertexAttribute(gpuObject, mesh.LightmapCoordinates(), 2)
		gosigl.FinishMesh()
		vals[idx] = gpuObject
	}
	ModelIdMap[msg.Resource.FilePath()] = vals
}

func DestroyPropOnGPU(dispatched event.IMessage) {
	msg := dispatched.(*message.PropUnloaded)
	for _, i := range ModelIdMap[msg.Resource.FilePath()] {
		gosigl.DeleteMesh(i)
	}
}
