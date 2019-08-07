package texture

import (
	"github.com/galaco/lambda-client/core/lib/math/shape"
	"github.com/galaco/packrect"
	"github.com/galaco/vtf/format"
	"github.com/go-gl/mathgl/mgl32"
)

// @TODO THIS IS NOT COMPLETE. IT DOES NOT WORK

// Atlas
// A texture atlas implementation.
type Atlas struct {
	Colour2D
}

// Format returns colour format
// For now always RGBA
func (atlas *Atlas) Format() uint32 {
	return uint32(format.RGB888)
}

// PackTextures
func (atlas *Atlas) PackTextures(textures []ITexture, padding int) ([]shape.Rect, error) {
	root := packrect.SubRect{
		Width:  atlas.Width(),
		Height: atlas.Height(),
	}
	rects := make([]packrect.IRectangle, len(textures))
	for idx, tex := range textures {
		rects[idx] = packrect.NewRectangle(tex.Width()+(padding*2), tex.Height()+(padding*2))
	}

	mapping, err := packrect.Pack(&root, rects)
	if err != nil {
		return nil, err
	}

	uvRects := make([]shape.Rect, len(mapping))
	for idx, rect := range mapping {
		atlas.writeTexture(textures[idx], &rect, padding)
		uvRects[idx] = *shape.NewRect(mgl32.Vec2{
			float32(rect.Left+padding) / float32(atlas.Width()),
			float32(rect.Top+padding) / float32(atlas.Height()),
		},
			mgl32.Vec2{
				float32(rect.Left+(rect.Width-(2*padding))) / float32(atlas.Width()),
				float32(rect.Top+(rect.Height-(2*padding))) / float32(atlas.Height()),
			})
	}

	return uvRects, nil
}

//// findSpace finds free space in atlas buffer to write rectangle to
//func (atlas *Atlas) findSpace(width int, height int) (x, y int, err error) {
//
//	return x, y, err
//}

// writeTexture insert write raw data into atlas at calculated position
func (atlas *Atlas) writeTexture(tex ITexture, location *packrect.SubRect, padding int) {
	data := tex.PixelDataForFrame(0)

	bytesPerPixel := 3
	bytesPerRow := bytesPerPixel * tex.Width()

	// indent into atlas top-left
	offset := (location.Top * (bytesPerPixel * atlas.Width())) + (bytesPerPixel * location.Left)

	// WRITE TOP PADDING
	for i := 0; i < padding; i++ {
		localOffset := 0
		row := make([]byte, bytesPerRow+2*(bytesPerPixel*padding))
		// left
		for j := 0; j < padding; j++ {
			copy(row[localOffset:localOffset+bytesPerPixel], data[(j*bytesPerPixel):(j*bytesPerPixel)+bytesPerPixel])
			localOffset += bytesPerPixel
		}
		// middle
		copy(row[localOffset:localOffset+bytesPerRow], data[:bytesPerRow])
		localOffset += bytesPerRow

		//right
		for j := 0; j < padding; j++ {
			copy(row[localOffset:localOffset+bytesPerPixel], data[(j*bytesPerPixel):(j*bytesPerPixel)+bytesPerPixel])
			localOffset += bytesPerPixel
		}

		// write row
		copy(atlas.rawColourData[offset:offset+len(row)], row[:])
		offset += (bytesPerPixel * atlas.Width())
	}

	// WRITE MAIN TEXTURE
	for i := 0; i < tex.Height(); i++ {
		localOffset := 0
		row := make([]byte, bytesPerRow+2*(bytesPerPixel*padding))
		// left
		for j := 0; j < padding; j++ {
			copy(row[localOffset:localOffset+bytesPerPixel], data[(j*bytesPerPixel):(j*bytesPerPixel)+bytesPerPixel])
			localOffset += bytesPerPixel
		}
		// middle
		copy(row[localOffset:localOffset+bytesPerRow], data[i*bytesPerRow:(i*bytesPerRow)+bytesPerRow])
		localOffset += bytesPerRow

		//right
		for j := 0; j < padding; j++ {
			copy(row[localOffset:localOffset+bytesPerPixel], data[(j*bytesPerPixel):(j*bytesPerPixel)+bytesPerPixel])
			localOffset += bytesPerPixel
		}

		// write row
		copy(atlas.rawColourData[offset:offset+len(row)], row[:])
		offset += (bytesPerPixel * atlas.Width())
	}

	// write bottom padding
	for i := 0; i < padding; i++ {
		localOffset := 0
		row := make([]byte, bytesPerRow+2*(bytesPerPixel*padding))
		// left
		for j := 0; j < padding; j++ {
			copy(row[localOffset:localOffset+bytesPerPixel], data[(j*bytesPerPixel):(j*bytesPerPixel)+bytesPerPixel])
			localOffset += bytesPerPixel
		}
		// middle
		copy(row[localOffset:localOffset+bytesPerRow], data[:bytesPerRow])
		localOffset += bytesPerRow

		//right
		for j := 0; j < padding; j++ {
			copy(row[localOffset:localOffset+bytesPerPixel], data[(j*bytesPerPixel):(j*bytesPerPixel)+bytesPerPixel])
			localOffset += bytesPerPixel
		}

		// write row
		copy(atlas.rawColourData[offset:offset+len(row)], row[:])
		offset += (bytesPerPixel * atlas.Width())
	}
}

// NewAtlas
func NewAtlas(width int, height int) *Atlas {
	return &Atlas{
		Colour2D: Colour2D{
			rawColourData: make([]uint8, width*height*3),
			Texture2D: Texture2D{
				width:  width,
				height: height,
			},
		},
	}
}
