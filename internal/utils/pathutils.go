package utils

import (
	"path/filepath"
	"strings"
)

// GetPathName extracts the filename without extension from a given path
func GetPathName(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// ChangePathExt changes the extension of the given path to the provided extension
func ChangePathExt(path string, ext string) string {
	dir := filepath.Dir(path)
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return filepath.Join(dir, base+ext)
}
