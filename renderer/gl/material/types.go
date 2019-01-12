package material

import "github.com/galaco/gosigl"

// getGLTextureFormat swap vtf format to openGL format
func GLTextureFormatFromVtfFormat(vtfFormat uint32) gosigl.PixelFormat {
	switch vtfFormat {
	case 0:
		return gosigl.RGBA
	case 2:
		return gosigl.RGB
	case 3:
		return gosigl.BGR
	case 12:
		return gosigl.BGRA
	case 13:
		return gosigl.DXT1
	case 14:
		return gosigl.DXT3
	case 15:
		return gosigl.DXT5
	default:
		return gosigl.RGB
	}
}

