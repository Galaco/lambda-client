package filesystem

import "strings"

// NormalisePath ensures that the same filepath format is used for paths,
// regardless of platform.
func NormalisePath(filePath string) string {
	return strings.Replace(filePath, "\\", "/", -1)
}
