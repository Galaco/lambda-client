package filesystem

import (
	"github.com/galaco/KeyValues"
	"github.com/galaco/lambda-client/core/lib/util"
	"github.com/galaco/lambda-client/core/lib/vpk"
	"path/filepath"
	"regexp"
	"strings"
)

// CreateFilesystemFromGameInfoDefinitions Reads game resource data paths
// from gameinfo.txt
// All games should ship with a gameinfo.txt, but it isn't actually mandatory.
func CreateFilesystemFromGameInfoDefinitions(basePath string, gameInfo *keyvalues.KeyValue) IFileSystem {
	fs := NewFileSystem()
	gameInfoNode, _ := gameInfo.Find("GameInfo")
	fsNode, _ := gameInfoNode.Find("FileSystem")

	searchPathsNode, _ := fsNode.Find("SearchPaths")
	searchPaths, _ := searchPathsNode.Children()
	basePath, _ = filepath.Abs(basePath)
	basePath = strings.Replace(basePath, "\\", "/", -1)

	for _, searchPath := range searchPaths {
		kv := searchPath
		path, _ := kv.AsString()
		path = strings.Trim(path, " ")

		// Current directory
		gameInfoPathRegex := regexp.MustCompile(`(?i)\|gameinfo_path\|`)
		if gameInfoPathRegex.MatchString(path) {
			path = gameInfoPathRegex.ReplaceAllString(path, basePath+"/")
		}

		// Executable directory
		allSourceEnginePathsRegex := regexp.MustCompile(`(?i)\|all_source_engine_paths\|`)
		if allSourceEnginePathsRegex.MatchString(path) {
			path = allSourceEnginePathsRegex.ReplaceAllString(path, basePath+"/../")
		}
		if strings.Contains(strings.ToLower(kv.Key()), "mod") && !strings.HasPrefix(path, basePath) {
			path = basePath + "/../" + path
		}

		// Strip vpk extension, then load it
		path = strings.Trim(strings.Trim(path, " "), "\"")
		if strings.HasSuffix(path, ".vpk") {
			path = strings.Replace(path, ".vpk", "", 1)
			vpkHandle, err := vpk.OpenVPK(path)
			if err != nil {
				util.Logger().Error(err)
				continue
			}
			fs.RegisterVpk(path, vpkHandle)
			util.Logger().Notice("Registered vpk: " + path)
		} else {
			// wildcard suffixes not useful
			if strings.HasSuffix(path, "/*") {
				path = strings.Replace(path, "/*", "", -1)
			}
			fs.RegisterLocalDirectory(path)
			util.Logger().Notice("Registered path: " + path)
		}
	}

	return fs
}
