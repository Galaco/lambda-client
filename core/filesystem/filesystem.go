package filesystem

import (
	"bytes"
	"github.com/galaco/bsp/lumps"
	"github.com/galaco/vpk2"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Path
type Path string

// FileSystem
type FileSystem struct {
	gameVPKs         map[Path]vpk.VPK
	localDirectories []string
	pakFile          *lumps.Pakfile
}

// NewFileSystem returns a new filesystem
func NewFileSystem() *FileSystem {
	return &FileSystem{
		gameVPKs:         map[Path]vpk.VPK{},
		localDirectories: make([]string, 0),
		pakFile:          nil,
	}
}

// PakFile returns loaded pakfile
// There can only be 1 registered pakfile at once.
func (fs *FileSystem) PakFile() *lumps.Pakfile {
	return fs.pakFile
}

// RegisterVpk registers a vpk package as a valid
// asset directory
func (fs *FileSystem) RegisterVpk(path string, vpkFile *vpk.VPK) {
	fs.gameVPKs[Path(path)] = *vpkFile
}

func (fs *FileSystem) UnregisterVpk(path string) {
	for key := range fs.gameVPKs {
		if key == Path(path) {
			delete(fs.gameVPKs, key)
		}
	}
}

// RegisterLocalDirectory register a filesystem path as a valid
// asset directory
func (fs *FileSystem) RegisterLocalDirectory(directory string) {
	fs.localDirectories = append(fs.localDirectories, directory)
}

func (fs *FileSystem) UnregisterLocalDirectory(directory string) {
	for idx, dir := range fs.localDirectories {
		if dir == directory {
			if len(fs.localDirectories) == 1 {
				fs.localDirectories = make([]string, 0)
				return
			}
			fs.localDirectories = append(fs.localDirectories[:idx], fs.localDirectories[idx+1:]...)
		}
	}
}

// RegisterPakFile Set a pakfile to be used as an asset directory.
// This would normally be called during each map load
func (fs *FileSystem) RegisterPakFile(pakFile *lumps.Pakfile) {
	fs.pakFile = pakFile
}

// UnregisterPakFile removes the current pakfile from
// available search locations
func (fs *FileSystem) UnregisterPakFile() {
	fs.pakFile = nil
}

// EnumerateResourcePaths returns all registered resource paths.
// PakFile is excluded.
func (fs *FileSystem) EnumerateResourcePaths() []string {
	list := make([]string, 0)

	for idx := range fs.gameVPKs {
		list = append(list, string(idx))
	}

	list = append(list, fs.localDirectories...)

	return list
}

// GetFile attempts to get stream for filename.
// Search order is Pak->FileSystem->VPK
func (fs *FileSystem) GetFile(filename string) (io.Reader, error) {
	// sanitise file
	searchPath := NormalisePath(strings.ToLower(filename))

	// try to read from pakfile first
	if fs.pakFile != nil {
		f, err := fs.pakFile.GetFile(searchPath)
		if err == nil && f != nil && len(f) != 0 {
			return bytes.NewReader(f), nil
		}
	}

	// Fallback to local filesystem
	for _, dir := range fs.localDirectories {
		if _, err := os.Stat(dir + "\\" + searchPath); os.IsNotExist(err) {
			continue
		}
		file, err := ioutil.ReadFile(dir + searchPath)
		if err != nil {
			return nil, err
		}
		return bytes.NewBuffer(file), nil
	}

	// Fall back to game vpk
	for _, vfs := range fs.gameVPKs {
		entry := vfs.Entry(searchPath)
		if entry != nil {
			return entry.Open()
		}
	}

	return nil, NewFileNotFoundError(filename)
}
