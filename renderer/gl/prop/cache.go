package prop

import (
	"github.com/galaco/Lambda-Core/core/event"
	"github.com/galaco/Lambda-Core/core/resource/message"
	"github.com/galaco/gosigl"
)

var ModelIdMap map[string][]*gosigl.VertexObject

func SyncPropToGpu(dispatched event.IMessage) {
	msg := dispatched.(*message.PropLoaded)
	if ModelIdMap[msg.Resource.GetFilePath()] != nil {
		return
	}
	vals := make([]*gosigl.VertexObject, len(msg.Resource.GetMeshes()))

	for idx, mesh := range msg.Resource.GetMeshes() {
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
	ModelIdMap[msg.Resource.GetFilePath()] = vals
}

func DestroyPropOnGPU(dispatched event.IMessage) {
	msg := dispatched.(*message.PropUnloaded)
	for _, i := range ModelIdMap[msg.Resource.GetFilePath()] {
		gosigl.DeleteMesh(i)
	}
}
