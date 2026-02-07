package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitRunner defines git operations
type GitRunner interface {
	Pull() error
	Push() error
	SetOrigin(url string, force bool) error
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
	return g.run("pull main main")
}

func (g RealGitRunner) Push() error {
	log.Printf("Pushing to %s\n", configFolder().TildePath)
	return g.run("push")
}

func (g RealGitRunner) SetOrigin(url string, force bool) error {
	// Check if repo is public (only for SSH URLs that can be converted to HTTPS)
	if !force {
		if httpsURL := sshToHTTPS(url); httpsURL != "" {
			if isPublicRepo(httpsURL) {
				return fmt.Errorf("repository appears to be public (accessible without authentication).\nUse --force to add this origin if you're sure")
			}
		}
	}

	// Auto-init if needed
	if _, err := os.Stat(filepath.Join(g.dir, ".git")); err != nil {
		if err := g.Init(); err != nil {
			return fmt.Errorf("git init failed: %w", err)
		}
	}

	// Set the remote
	return g.run("remote", "add", "origin", url)
}

// sshToHTTPS converts git@github.com:user/repo.git to https://github.com/user/repo
// Returns empty string if not an SSH URL
func sshToHTTPS(url string) string {
	if strings.HasPrefix(url, "git@") {
		// git@github.com:user/repo.git -> https://github.com/user/repo
		parts := strings.TrimPrefix(url, "git@")
		colonIdx := strings.Index(parts, ":")
		if colonIdx > 0 {
			host := parts[:colonIdx]
			repo := strings.TrimSuffix(parts[colonIdx+1:], ".git")
			return fmt.Sprintf("https://%s/%s", host, repo)
		}
	}
	return ""
}

// isPublicRepo checks if a repo is publicly accessible via HTTP
func isPublicRepo(httpsURL string) bool {
	resp, err := http.Head(httpsURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
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
