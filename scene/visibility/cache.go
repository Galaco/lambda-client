package visibility

type Cache struct {
	Leafs      []uint16
	ClusterId  int16
	SkyVisible bool
	Faces      []uint16
}
