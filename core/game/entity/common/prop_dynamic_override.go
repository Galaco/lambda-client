package common

import (
	"github.com/galaco/lambda-client/core/entity"
	entity2 "github.com/galaco/lambda-client/core/game/entity"
)

// PropDynamicOverride
type PropDynamicOverride struct {
	entity.Base
	entity2.PropBase
}

// New
func (entity *PropDynamicOverride) New() entity.IEntity {
	return &PropDynamicOverride{}
}

// Classname
func (entity PropDynamicOverride) Classname() string {
	return "prop_dynamic_override"
}
