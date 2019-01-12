package material

import (
	"github.com/galaco/Lambda-Core/core/event"
	"github.com/galaco/Lambda-Core/core/resource"
	"github.com/galaco/Lambda-Core/core/resource/message"
	"github.com/galaco/gosigl"
)

var TextureIdMap map[string]gosigl.TextureBindingId

func SyncTextureToGpu(dispatched event.IMessage) {
	msg := dispatched.(*message.TextureLoaded)
	TextureIdMap[msg.Resource.(resource.IResource).GetFilePath()] = gosigl.CreateTexture2D(
		gosigl.TextureSlot(0),
		msg.Resource.Width(),
		msg.Resource.Height(),
		msg.Resource.PixelDataForFrame(0),
		gosigl.PixelFormat(GLTextureFormatFromVtfFormat(msg.Resource.Format())),
		false)
}

func DestroyTextureOnGPU(dispatched event.IMessage) {
	msg := dispatched.(*message.TextureUnloaded)
	gosigl.DeleteTextures(TextureIdMap[msg.Resource.GetFilePath()])
}
