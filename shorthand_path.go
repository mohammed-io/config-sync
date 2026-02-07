package main

import (
	"log"
	"os/user"
	"path/filepath"
	"strings"
)

// ShorthandPath represents a file path with both full and tilde (~) notations
type ShorthandPath struct {
	FullPath  string
	TildePath string
}

// Suffix returns a new ShorthandPath with the given suffix appended to both paths
func (ShorthandPath) New(str string) ShorthandPath {
	return ShorthandPath{
		TildePath: collapseToTilde(str),
		FullPath:  expandFromTilde(str),
	}
}

// Suffix returns a new ShorthandPath with the given suffix appended to both paths
func (self ShorthandPath) Suffix(str string) ShorthandPath {
	return ShorthandPath{
		TildePath: filepath.Join(self.TildePath, str),
		FullPath:  filepath.Join(self.FullPath, str),
	}
}

// expandFromTilde converts a tilde-prefixed path to an absolute path
func expandFromTilde(path string) string {
	if !strings.HasPrefix(path, "~/") && path != "~" {
		absolutePath, _ := filepath.Abs(path)
		return absolutePath
	}

	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("Could not get the current user information")
	}

	homeDir := currentUser.HomeDir

	if path == "~" {
		return homeDir
	}

	return filepath.Clean(filepath.Join(homeDir, path[2:]))
}

// collapseToTilde converts an absolute home-dir path to tilde notation
func collapseToTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		return path
	}

	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("Could not get the current user information")
	}
	homeDir := currentUser.HomeDir

	return filepath.Clean(filepath.Join("~/", path[len(homeDir):]))
}
