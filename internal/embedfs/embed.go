package embedfs

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed frontend/*
var frontendFS embed.FS

// GetFrontendFS returns the embedded React frontend filesystem.
// Returns nil if the frontend hasn't been built yet.
func GetFrontendFS() fs.FS {
	sub, err := fs.Sub(frontendFS, "frontend")
	if err != nil {
		return nil
	}
	return sub
}

// HasFrontend returns true if the frontend has been built and embedded.
func HasFrontend() bool {
	return GetFrontendFS() != nil
}

// PlaceholderDirExists checks if the dist directory exists on disk (for development).
func PlaceholderDirExists() bool {
	info, err := os.Stat(filepath.Join("ui", "dist"))
	if err == nil && info.IsDir() {
		return true
	}
	return false
}