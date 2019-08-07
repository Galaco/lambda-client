package scene

import (
	"github.com/galaco/lambda-client/core/model"
)

// IScene
type IScene interface {
	// Bsp
	Bsp() *model.Bsp
	// StaticProps
	StaticProps() []model.StaticProp
}
