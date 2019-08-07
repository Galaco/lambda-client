package entity

import (
	"github.com/galaco/lambda-client/core/entity"
	"github.com/galaco/lambda-client/core/filesystem"
	entity2 "github.com/galaco/lambda-client/core/game/entity"
	"github.com/galaco/lambda-client/core/loader/prop"
	"github.com/galaco/lambda-client/core/resource"
	"strings"
)

// DoesEntityReferenceStudioModel tests if an entity is
// tied to a model (normally prop_* classnames, but not exclusively)
func DoesEntityReferenceStudioModel(ent entity.IEntity) bool {
	return strings.HasSuffix(ent.KeyValues().ValueForKey("model"), ".mdl")
}

// AssignStudioModelToEntity sets a renderable entity's model
func AssignStudioModelToEntity(entity entity.IEntity, fs filesystem.IFileSystem) {
	modelName := entity.KeyValues().ValueForKey("model")
	if !resource.Manager().HasModel(modelName) {
		m, _ := prop.LoadProp(modelName, fs)
		entity.(entity2.IProp).SetModel(m)
	} else {
		entity.(entity2.IProp).SetModel(resource.Manager().Model(modelName))
	}
}
