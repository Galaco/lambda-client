package common

import (
	entity2 "github.com/galaco/lambda-client/core/entity"
	"github.com/galaco/lambda-client/core/game/entity"
)

// PropDynamic
type PropDynamic struct {
	entity2.Base
	entity.PropBase
}

// New
func (entity *PropDynamic) New() entity2.IEntity {
	return &PropDynamic{}
}

// Classname
func (entity PropDynamic) Classname() string {
	return "prop_dynamic"
}
