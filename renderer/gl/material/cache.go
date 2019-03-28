package material

import (
	"github.com/galaco/Lambda-Core/core/event"
	"github.com/galaco/Lambda-Core/core/material"
	"github.com/galaco/Lambda-Core/core/resource/message"
	"github.com/galaco/gosigl"
)

type Cache struct {
	textureIdMap map[string]gosigl.TextureBindingId
}

func (cache *Cache) FetchCachedTexture(textureName string) gosigl.TextureBindingId {
	return cache.textureIdMap[textureName]
}

func (cache *Cache) SyncTextureToGpu(dispatched event.IMessage) {
	msg := dispatched.(*message.MaterialLoaded)
	mat := msg.Resource.(*material.Material)

	if mat.Textures.Albedo == nil {
		return
	}

	if _,ok := cache.textureIdMap[mat.Textures.Albedo.GetFilePath()]; !ok {
		cache.textureIdMap[mat.Textures.Albedo.GetFilePath()] = gosigl.CreateTexture2D(
			gosigl.TextureSlot(0),
			mat.Textures.Albedo.Width(),
			mat.Textures.Albedo.Height(),
			mat.Textures.Albedo.PixelDataForFrame(0),
			gosigl.PixelFormat(GLTextureFormatFromVtfFormat(mat.Textures.Albedo.Format())),
			false)
	}

	if mat.Textures.Normal != nil {
		if _,ok := cache.textureIdMap[mat.Textures.Normal.GetFilePath()]; !ok {
			cache.textureIdMap[mat.Textures.Normal.GetFilePath()] = gosigl.CreateTexture2D(
				gosigl.TextureSlot(1),
				mat.Textures.Normal.Width(),
				mat.Textures.Normal.Height(),
				mat.Textures.Normal.PixelDataForFrame(0),
				gosigl.PixelFormat(GLTextureFormatFromVtfFormat(mat.Textures.Normal.Format())),
				false)
		}
	}
}

func (cache *Cache) DestroyTextureOnGPU(dispatched event.IMessage) {
	msg := dispatched.(*message.MaterialLoaded)
	mat := msg.Resource.(*material.Material)
	if mat.Textures.Albedo != nil {
		gosigl.DeleteTextures(cache.textureIdMap[mat.Textures.Albedo.GetFilePath()])
	}
	if mat.Textures.Normal != nil {
		gosigl.DeleteTextures(cache.textureIdMap[mat.Textures.Normal.GetFilePath()])
	}
}

func NewCache() *Cache {
	return &Cache{
		textureIdMap: map[string]gosigl.TextureBindingId{},
	}
}
