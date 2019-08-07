package filesystem

import (
	"github.com/galaco/bsp/lumps"
	vpk "github.com/galaco/vpk2"
	"io"
)

// IFileSystem represents a Source Engine filesystem
type IFileSystem interface {
	// PakFile returns a loaded map pakfile
	PakFile() *lumps.Pakfile
	// RegisterPakFile adds a maps pakfile.
	RegisterPakFile(pakfile *lumps.Pakfile)
	// RegisterVpk adds a VPK to this filesystem
	RegisterVpk(path string, vpkFile *vpk.VPK)
	// RegisterLocalDirectory adds a local directory
	RegisterLocalDirectory(directory string)
	// UnregisterLocalDirectory removes a registered local directory
	UnregisterLocalDirectory(directory string)
	// UnregisterPakFile removed a loaded pakfile
	UnregisterPakFile()
	// EnumerateResourcePaths returns all loaded paths (local, vpk locations). Does not include pakfile.
	EnumerateResourcePaths() []string
	// GetFile returns a file, or error if not found
	GetFile(filename string) (io.Reader, error)
}
