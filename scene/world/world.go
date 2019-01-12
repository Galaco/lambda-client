package world

import (
	"github.com/galaco/Lambda-Client/scene/visibility"
	"github.com/galaco/Lambda-Core/core/entity"
	"github.com/galaco/Lambda-Core/core/mesh"
	"github.com/galaco/Lambda-Core/core/model"
	"github.com/galaco/bsp/primitives/leaf"
	"github.com/go-gl/mathgl/mgl32"
	"sync"
)

type World struct {
	entity.Base
	staticModel model.Bsp

	staticProps []model.StaticProp
	sky         Sky

	visibleClusterLeafs []*model.ClusterLeaf

	visData     *visibility.Vis
	LeafCache   *visibility.Cache
	currentLeaf *leaf.Leaf

	rebuildMutex sync.Mutex
}

func (entity *World) Bsp() *model.Bsp {
	return &entity.staticModel
}

func (entity *World) Sky() *Sky {
	return &entity.sky
}

func (entity *World) VisibleClusters() []*model.ClusterLeaf {
	entity.rebuildMutex.Lock()
	vw := entity.visibleClusterLeafs
	entity.rebuildMutex.Unlock()
	return vw
}

// Rebuild the current facelist to render, by first
// recalculating using vvis data
func (entity *World) TestVisibility(position mgl32.Vec3) {
	// View hasn't moved
	currentLeaf := entity.visData.FindCurrentLeaf(position)

	if currentLeaf == entity.currentLeaf {
		return
	}

	if currentLeaf == nil || currentLeaf.Cluster == -1 {
		// Still outside the world
		if entity.currentLeaf == nil {
			return
		}

		entity.currentLeaf = currentLeaf

		entity.AsyncRebuildVisibleWorld()
		return
	}

	// Haven't changed cluster
	if entity.LeafCache != nil && entity.LeafCache.ClusterId == currentLeaf.Cluster {
		return
	}

	entity.currentLeaf = currentLeaf
	entity.LeafCache = entity.visData.GetPVSCacheForCluster(currentLeaf.Cluster)

	entity.AsyncRebuildVisibleWorld()
}

// Launches rebuilding the visible world in a separate thread
// Note: This *could* cause rendering issues if the rebuild is slower than
// travelling between clusters
func (entity *World) AsyncRebuildVisibleWorld() {
	func(currentLeaf *leaf.Leaf) {
		visibleWorld := make([]*model.ClusterLeaf, 0)

		visibleClusterIds := make([]int16, 0)

		if currentLeaf != nil && currentLeaf.Cluster != -1 {
			visibleClusterIds = entity.visData.PVSForCluster(currentLeaf.Cluster)
		}

		// nothing visible so render everything
		if len(visibleClusterIds) == 0 {
			for idx := range entity.staticModel.ClusterLeafs() {
				visibleWorld = append(visibleWorld, &entity.staticModel.ClusterLeafs()[idx])
			}
		} else {
			for _, clusterId := range visibleClusterIds {
				visibleWorld = append(visibleWorld, &entity.staticModel.ClusterLeafs()[clusterId])
			}
		}

		entity.rebuildMutex.Lock()
		entity.visibleClusterLeafs = visibleWorld
		entity.rebuildMutex.Unlock()
	}(entity.currentLeaf)
}

// Build skybox from tree
func (entity *World) BuildSkybox(sky *model.Model, position mgl32.Vec3, scale float32) {
	// Rebuild bsp faces
	visibleModel := model.NewBsp(entity.staticModel.Mesh().(*mesh.Mesh))

	visibleWorld := make([]*model.ClusterLeaf, 0)

	l := entity.visData.FindCurrentLeaf(position)
	visibleClusterIds := entity.visData.PVSForCluster(l.Cluster)

	// nothing visible so render everything
	if len(visibleClusterIds) == 0 {
		for idx := range entity.staticModel.ClusterLeafs() {
			visibleWorld = append(visibleWorld, &entity.staticModel.ClusterLeafs()[idx])
		}
	} else {
		for clusterId := range visibleClusterIds {
			visibleWorld = append(visibleWorld, &entity.staticModel.ClusterLeafs()[clusterId])
		}
	}

	entity.sky = *NewSky(visibleModel, visibleWorld, position, scale, sky)
}

func NewWorld(world model.Bsp, staticProps []model.StaticProp, visData *visibility.Vis) *World {
	c := World{
		staticModel: world,
		staticProps: staticProps,
		visData:     visData,
	}

	c.TestVisibility(mgl32.Vec3{0, 0, 0})

	return &c
}
