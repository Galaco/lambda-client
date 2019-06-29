package world

import (
	"github.com/galaco/lambda-core/entity"
	"github.com/galaco/lambda-core/model"
	"github.com/go-gl/mathgl/mgl32"
)

type Sky struct {
	geometry     *model.Bsp
	clusterLeafs []*model.ClusterLeaf
	transform    entity.Transform

	cubemap *model.Model
}

func (sky *Sky) GetVisibleBsp() *model.Bsp {
	return sky.geometry
}

func (sky *Sky) GetClusterLeafs() []*model.ClusterLeaf {
	return sky.clusterLeafs
}

func (sky *Sky) GetCubemap() *model.Model {
	return sky.cubemap
}

func (sky *Sky) Transform() *entity.Transform {
	return &sky.transform
}

func NewSky(bsp *model.Bsp, clusterLeafs []*model.ClusterLeaf, position mgl32.Vec3, scale float32, skyCube *model.Model) *Sky {
	s := Sky{
		geometry:     bsp,
		clusterLeafs: clusterLeafs,
		cubemap:      skyCube,
	}

	skyCameraPosition := (mgl32.Vec3{0, 0, 0}).Sub(position)
	skyCameraScale := mgl32.Vec3{scale, scale, scale}

	s.transform.Position = skyCameraPosition.Mul(scale)
	s.transform.Scale = skyCameraScale

	// remap prop transform to real world
	for _, l := range s.clusterLeafs {
		for _, prop := range l.StaticProps {
			prop.Transform().Position = prop.Transform().Position.Add(skyCameraPosition)
			prop.Transform().Position = prop.Transform().Position.Mul(scale)
			prop.Transform().Scale = skyCameraScale
		}
	}
	return &s
}
