package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// GitRunner defines git operations
type GitRunner interface {
	Pull() error
	Push() error
	SetOrigin(url string) error
	Add() error
	Commit(message string) error
	AddAndPush(message string) error
	Init() error
}

// RealGitRunner executes actual git commands
type RealGitRunner struct {
	dir string
}

func (g RealGitRunner) run(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (g RealGitRunner) Pull() error {
	log.Printf("Pulling from %s\n", configFolder().TildePath)
	return g.run("pull")
}

func (g RealGitRunner) Push() error {
	log.Printf("Pushing to %s\n", configFolder().TildePath)
	return g.run("push")
}

func (g RealGitRunner) SetOrigin(url string) error {
	// Auto-init if needed
	if _, err := os.Stat(filepath.Join(g.dir, ".git")); err != nil {
		if err := g.Init(); err != nil {
			return fmt.Errorf("git init failed: %w", err)
		}
	}

	// Set the remote
	return g.run("remote", "add", "origin", url)
}

func (g RealGitRunner) Add() error {
	return g.run("add", "-A")
}

func (g RealGitRunner) Commit(message string) error {
	return g.run("commit", "-m", message)
}

func (g RealGitRunner) AddAndPush(message string) error {
	if err := g.Add(); err != nil {
		return err
	}
	return g.Commit(message)
}

func (g RealGitRunner) Init() error {
	// Check if .git exists
	if _, err := os.Stat(filepath.Join(g.dir, ".git")); err == nil {
		return nil // Already initialized
	}
	log.Printf("Initializing git repository in %s\n", configFolder().TildePath)
	return g.run("init")
}

// NewGitRunner creates a new GitRunner for the config folder
func NewGitRunner() GitRunner {
	return RealGitRunner{dir: configFolder().FullPath}
}
