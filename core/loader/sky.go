package loader

import (
	"github.com/galaco/lambda-client/core/filesystem"
	material2 "github.com/galaco/lambda-client/core/loader/material"
	"github.com/galaco/lambda-client/core/material"
	"github.com/galaco/lambda-client/core/mesh/primitive"
	"github.com/galaco/lambda-client/core/model"
	"github.com/galaco/lambda-client/core/texture"
)

const skyboxRootDir = "skybox/"

// LoadSky loads the skymaterial cubemap.
// The materialname is normally obtained from the worldspawn entity
func LoadSky(materialName string, fs filesystem.IFileSystem) *model.Model {
	sky := model.NewModel(materialName)

	mats := make([]material.IMaterial, 6)

	mats[0] = material2.LoadSingleMaterial(skyboxRootDir+materialName+"up.vmt", fs)
	mats[1] = material2.LoadSingleMaterial(skyboxRootDir+materialName+"dn.vmt", fs)
	mats[2] = material2.LoadSingleMaterial(skyboxRootDir+materialName+"lf.vmt", fs)
	mats[3] = material2.LoadSingleMaterial(skyboxRootDir+materialName+"rt.vmt", fs)
	mats[4] = material2.LoadSingleMaterial(skyboxRootDir+materialName+"ft.vmt", fs)
	mats[5] = material2.LoadSingleMaterial(skyboxRootDir+materialName+"bk.vmt", fs)

	texs := make([]texture.ITexture, 6)
	for i := 0; i < 6; i++ {
		texs[i] = mats[i].(*material.Material).Textures.Albedo
	}

	sky.AddMesh(primitive.NewCube())

	sky.Meshes()[0].SetMaterial(texture.NewCubemap(texs))

	return sky
}
