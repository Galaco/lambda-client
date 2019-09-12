package material

import (
	"github.com/galaco/bsp/lumps"
	"github.com/galaco/bsp/primitives/texinfo"
	"github.com/galaco/lambda-client/core/filesystem"
	"github.com/galaco/lambda-client/core/lib/stringtable"
	"github.com/galaco/source-tools-common/texdatastringtable"
)

// LoadMaterials is the base bsp material loader function.
// All bsp materials should be loaded by this function.
// Note that this covers bsp referenced materials only, model & entity
// materials are loaded mostly ad-hoc.
func LoadMaterials(fs filesystem.IFileSystem, stringData *lumps.TexDataStringData, stringTable *lumps.TexDataStringTable, texInfos *[]texinfo.TexInfo) *texdatastringtable.TexDataStringTable {
	materialStringTable := stringtable.NewTable(stringData, stringTable)
	LoadErrorMaterial()
	LoadMaterialList(fs, stringtable.SortUnique(materialStringTable, texInfos))

	return materialStringTable
}
