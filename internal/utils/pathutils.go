package utils

import (
	"path/filepath"
	"strings"
)

// GetPathName returns the last part of the path
func GetPathName(path string) string {
	LogInfo("Getting path name from", path)
	return strings.Split(strings.Split(path, "/")[len(strings.Split(path, "/"))-1], ".")[0]
}


// ChangePathExt changes the extension of the given path to the provided ext.
func ChangePathExt(path string, ext string) string {
	// Get the directory and base name of the path
	dir := filepath.Dir(path)
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	
	// Create the new path with the desired extension
	out := filepath.Join(dir, base+ext)
	
	// Log the change
	LogInfo("Changing path extension from " + path + " to " + out)
	
	return out
}
