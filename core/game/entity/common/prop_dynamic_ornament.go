package common

import (
	"github.com/galaco/lambda-client/core/entity"
	entity2 "github.com/galaco/lambda-client/core/game/entity"
)

//PropDynamicOrnament
type PropDynamicOrnament struct {
	entity.Base
	entity2.PropBase
}

// New
func (entity *PropDynamicOrnament) New() entity.IEntity {
	return &PropDynamicOrnament{}
}

// Classname
func (entity PropDynamicOrnament) Classname() string {
	return "prop_dynamic_ornament"
}
