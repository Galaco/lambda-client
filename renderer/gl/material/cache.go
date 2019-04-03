package material

import (
	"github.com/galaco/Lambda-Core/core/event"
	"github.com/galaco/Lambda-Core/core/material"
	"github.com/galaco/Lambda-Core/core/resource/message"
	"github.com/galaco/gosigl"
	"sync"
)

type Cache struct {
	textureIdMap map[string]gosigl.TextureBindingId
	mut          sync.Mutex
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

	cache.mut.Lock()
	if _, ok := cache.textureIdMap[mat.Textures.Albedo.FilePath()]; !ok {
		cache.textureIdMap[mat.Textures.Albedo.FilePath()] = gosigl.CreateTexture2D(
			gosigl.TextureSlot(0),
			mat.Textures.Albedo.Width(),
			mat.Textures.Albedo.Height(),
			mat.Textures.Albedo.PixelDataForFrame(0),
			gosigl.PixelFormat(GLTextureFormatFromVtfFormat(mat.Textures.Albedo.Format())),
			false)
	}

	if mat.Textures.Normal != nil {
		if _, ok := cache.textureIdMap[mat.Textures.Normal.FilePath()]; !ok {
			cache.textureIdMap[mat.Textures.Normal.FilePath()] = gosigl.CreateTexture2D(
				gosigl.TextureSlot(1),
				mat.Textures.Normal.Width(),
				mat.Textures.Normal.Height(),
				mat.Textures.Normal.PixelDataForFrame(0),
				gosigl.PixelFormat(GLTextureFormatFromVtfFormat(mat.Textures.Normal.Format())),
				false)
		}
	}
	cache.mut.Unlock()
}

func (cache *Cache) DestroyTextureOnGPU(dispatched event.IMessage) {
	msg := dispatched.(*message.MaterialUnloaded)
	mat := msg.Resource.(*material.Material)
	if mat.Textures.Albedo != nil {
		gosigl.DeleteTextures(cache.textureIdMap[mat.Textures.Albedo.FilePath()])
	}
	if mat.Textures.Normal != nil {
		gosigl.DeleteTextures(cache.textureIdMap[mat.Textures.Normal.FilePath()])
	}
}

func NewCache() *Cache {
	return &Cache{
		textureIdMap: map[string]gosigl.TextureBindingId{},
	}
}
