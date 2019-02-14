package visibility

import (
	"github.com/galaco/bsp"
	"github.com/galaco/bsp/lumps"
	"github.com/galaco/bsp/primitives/leaf"
	"github.com/galaco/bsp/primitives/node"
	"github.com/galaco/bsp/primitives/plane"
	"github.com/galaco/bsp/primitives/visibility"
	"github.com/go-gl/mathgl/mgl32"
)

type Vis struct {
	ClusterCache   []Cache
	VisibilityLump *visibility.Vis
	Leafs          []leaf.Leaf
	LeafFaces      []uint16
	Nodes          []node.Node
	Planes         []plane.Plane

	viewPosition    mgl32.Vec3
	viewCurrentLeaf *leaf.Leaf
}

func (vis *Vis) PVSForCluster(clusterId int16) []int16 {
	return vis.VisibilityLump.GetVisibleClusters(clusterId)
}

func (vis *Vis) GetPVSCacheForCluster(clusterId int16) *Cache {
	if clusterId == -1 {
		clusterId = int16(vis.findCurrentLeafIndex(vis.viewPosition))
	}
	for _, cacheEntry := range vis.ClusterCache {
		if cacheEntry.ClusterId == clusterId {
			return &cacheEntry
		}
	}
	return vis.cachePVSForCluster(clusterId)
}

// Cache visible data for current cluster
func (vis *Vis) cachePVSForCluster(clusterId int16) *Cache {
	clusterList := vis.VisibilityLump.GetPVSForCluster(clusterId)

	skyVisible := false

	faces := make([]uint16, 0)
	leafs := make([]uint16, 0)
	for idx, l := range vis.Leafs {
		//Check if cluster is in pvs
		if !vis.clusterVisible(&clusterList, l.Cluster) {
			continue
		}
		if l.Flags()&leaf.LEAF_FLAGS_SKY > 0 {
			skyVisible = true
		}
		leafs = append(leafs, uint16(idx))
		faces = append(faces, vis.LeafFaces[l.FirstLeafFace:l.FirstLeafFace+l.NumLeafFaces]...)
	}

	cache := Cache{
		ClusterId:  clusterId,
		Faces:      faces,
		Leafs:      leafs,
		SkyVisible: skyVisible,
	}

	vis.ClusterCache = append(vis.ClusterCache, cache)

	return &cache
}

// Determine if a cluster is visible
func (vis *Vis) clusterVisible(pvs *[]bool, leafCluster int16) bool {
	if leafCluster < 0 {
		return true
	}

	if (*pvs)[leafCluster] {
		return true
	}

	return false
}

// Test if the camera has moved, and find the current leaf if so
func (vis *Vis) FindCurrentLeaf(position mgl32.Vec3) *leaf.Leaf {
	if !vis.viewPosition.ApproxEqualThreshold(position, 0.000000001) {
		vis.viewPosition = position
		vis.viewCurrentLeaf = &vis.Leafs[vis.findCurrentLeafIndex(vis.viewPosition)]
	}
	return vis.viewCurrentLeaf
}

// Find the index into the leaf array for the leaf the player
// is inside of
// Based on: https://bitbucket.org/fallahn/chuf-arc
func (vis *Vis) findCurrentLeafIndex(position mgl32.Vec3) int32 {
	i := int32(0)

	//walk the bsp to find the index of the leaf which contains our position
	for i >= 0 {
		node := &vis.Nodes[i]
		plane := vis.Planes[node.PlaneNum]

		//check which side of the plane the position is on so we know which direction to go
		distance := plane.Normal.X()*position.X() + plane.Normal.Y()*position.Y() + plane.Normal.Z()*position.Z() - plane.Distance
		i = node.Children[0]
		if distance < 0 {
			i = node.Children[1]
		}
	}

	return ^i
}

func NewVisFromBSP(file *bsp.Bsp) *Vis {
	return &Vis{
		VisibilityLump: file.GetLump(bsp.LUMP_VISIBILITY).(*lumps.Visibility).GetData(),
		viewPosition:   mgl32.Vec3{65536, 65536, 65536},
		Leafs:          file.GetLump(bsp.LUMP_LEAFS).(*lumps.Leaf).GetData(),
		LeafFaces:      file.GetLump(bsp.LUMP_LEAFFACES).(*lumps.LeafFace).GetData(),
		Nodes:          file.GetLump(bsp.LUMP_NODES).(*lumps.Node).GetData(),
		Planes:         file.GetLump(bsp.LUMP_PLANES).(*lumps.Planes).GetData(),
	}
}
