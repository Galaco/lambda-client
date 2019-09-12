package bsp

import (
	"github.com/galaco/gosigl"
	"github.com/galaco/lambda-client/core/event"
	"github.com/galaco/lambda-client/core/resource/message"
)

var MapGPUResource *gosigl.VertexObject

func SyncMapToGpu(dispatched event.IMessage) {
	msg := dispatched.(*message.MapLoaded)
	mesh := msg.Resource.Mesh()
	MapGPUResource = gosigl.NewMesh(mesh.Vertices())
	gosigl.CreateVertexAttribute(MapGPUResource, mesh.UVs(), 2)
	gosigl.CreateVertexAttribute(MapGPUResource, mesh.Normals(), 3)
	if len(mesh.Tangents()) == 0 {
		mesh.GenerateTangents()
	}
	gosigl.CreateVertexAttribute(MapGPUResource, mesh.Tangents(), 4)
	gosigl.CreateVertexAttribute(MapGPUResource, mesh.LightmapCoordinates(), 2)
	gosigl.FinishMesh()
}
