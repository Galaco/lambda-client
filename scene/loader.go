package scene

import (
	bsplib "github.com/galaco/bsp"
	"github.com/galaco/bsp/lumps"
	"github.com/galaco/lambda-client/internal/config"
	"github.com/galaco/lambda-client/scene/visibility"
	"github.com/galaco/lambda-client/scene/world"
	"github.com/galaco/lambda-core/entity"
	"github.com/galaco/lambda-core/lib/util"
	"github.com/galaco/lambda-core/loader"
	entity2 "github.com/galaco/lambda-core/loader/entity"
	"github.com/galaco/lambda-core/model"
	entitylib "github.com/galaco/source-tools-common/entity"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/filesystem"
)

func LoadFromFile(fileName string, fs *filesystem.FileSystem) {
	newScene := Get()

	bspData, err := bsplib.ReadFromFile(fileName)
	if err != nil {
		util.Logger().Panic(err)
	}
	if bspData.Header().Version < 19 {
		util.Logger().Panic("Unsupported BSP Version. Exiting...")
	}

	//Set pakfile for filesystem
	fs.RegisterPakFile(bspData.Lump(bsplib.LumpPakfile).(*lumps.Pakfile))

	loadWorld(newScene, bspData, fs)

	loadEntities(newScene, bspData.Lump(bsplib.LumpEntities).(*lumps.EntData), fs)

	loadCamera(newScene)
}

func loadWorld(targetScene *Scene, file *bsplib.Bsp, fs *filesystem.FileSystem) {
	baseWorld := loader.LoadMap(fs, file)

	baseWorldBsp := baseWorld.Bsp()
	baseWorldBspFaces := baseWorldBsp.ClusterLeafs()[0].Faces
	baseWorldStaticProps := baseWorld.StaticProps()

	visData := visibility.NewVisFromBSP(file)
	bspClusters := make([]model.ClusterLeaf, visData.VisibilityLump.NumClusters)
	defaultCluster := model.ClusterLeaf{
		Id: 32767,
	}
	for _, bspLeaf := range visData.Leafs {
		for _, leafFace := range visData.LeafFaces[bspLeaf.FirstLeafFace : bspLeaf.FirstLeafFace+bspLeaf.NumLeafFaces] {
			if bspLeaf.Cluster == -1 {
				//defaultCluster.Faces = append(defaultCluster.Faces, bspFaces[leafFace])
				continue
			}
			bspClusters[bspLeaf.Cluster].Id = bspLeaf.Cluster
			bspClusters[bspLeaf.Cluster].Faces = append(bspClusters[bspLeaf.Cluster].Faces, baseWorldBspFaces[leafFace])
			bspClusters[bspLeaf.Cluster].Mins = mgl32.Vec3{
				float32(bspLeaf.Mins[0]),
				float32(bspLeaf.Mins[1]),
				float32(bspLeaf.Mins[2]),
			}
			bspClusters[bspLeaf.Cluster].Maxs = mgl32.Vec3{
				float32(bspLeaf.Maxs[0]),
				float32(bspLeaf.Maxs[1]),
				float32(bspLeaf.Maxs[2]),
			}
			bspClusters[bspLeaf.Cluster].Origin = bspClusters[bspLeaf.Cluster].Mins.Add(bspClusters[bspLeaf.Cluster].Maxs.Sub(bspClusters[bspLeaf.Cluster].Mins))
		}
	}

	// Assign staticprops to clusters
	for idx, prop := range baseWorld.StaticProps() {
		for _, leafId := range prop.LeafList() {
			clusterId := visData.Leafs[leafId].Cluster
			if clusterId == -1 {
				defaultCluster.StaticProps = append(defaultCluster.StaticProps, &baseWorldStaticProps[idx])
				continue
			}
			bspClusters[clusterId].StaticProps = append(bspClusters[clusterId].StaticProps, &baseWorldStaticProps[idx])
		}
	}

	for _, idx := range baseWorldBsp.ClusterLeafs()[0].DispFaces {
		defaultCluster.Faces = append(defaultCluster.Faces, baseWorldBspFaces[idx])
	}

	baseWorldBsp.SetClusterLeafs(bspClusters)
	baseWorldBsp.SetDefaultCluster(defaultCluster)

	targetScene.SetWorld(world.NewWorld(*baseWorld.Bsp(), baseWorld.StaticProps(), visData))
}

func loadEntities(targetScene *Scene, entdata *lumps.EntData, fs *filesystem.FileSystem) {
	vmfEntityTree, err := entity2.ParseEntities(entdata.GetData())
	if err != nil {
		util.Logger().Panic(err)
	}
	entityList := entitylib.FromVmfNodeTree(vmfEntityTree.Unclassified)
	util.Logger().Notice("Found %d entities\n", entityList.Length())
	for i := 0; i < entityList.Length(); i++ {
		targetScene.AddEntity(entity2.CreateEntity(entityList.Get(i), fs))
	}

	skyCamera := entityList.FindByKeyValue("classname", "sky_camera")
	if skyCamera == nil {
		return
	}

	worldSpawn := entityList.FindByKeyValue("classname", "worldspawn")
	if worldSpawn == nil {
		return
	}

	targetScene.world.BuildSkybox(
		loader.LoadSky(worldSpawn.ValueForKey("skyname"), fs),
		skyCamera.VectorForKey("origin"),
		float32(skyCamera.IntForKey("scale")))
}

func loadCamera(targetScene *Scene) {
	targetScene.AddCamera(entity.NewCamera(mgl32.DegToRad(70), float32(config.Get().Video.Width)/float32(config.Get().Video.Height)))
}
