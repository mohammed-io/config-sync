package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// TimeoutConfig defines timeouts for different operations
type TimeoutConfig struct {
	GitRemote time.Duration // ls-remote (default 5s)
	GitLocal  time.Duration // status, rev-list, diff (default 2s)
	FileCopy  time.Duration // per file copy (default 1s)
	FileStat  time.Duration // file stat checks (default 500ms)
}

// DefaultTimeouts returns default timeout configuration
func DefaultTimeouts() TimeoutConfig {
	return TimeoutConfig{
		GitRemote: 5 * time.Second,
		GitLocal:  2 * time.Second,
		FileCopy:  1 * time.Second,
		FileStat:  500 * time.Millisecond,
	}
}

// TimingLogger logs operation durations and timeouts
type TimingLogger struct {
	prefix string
	enabled bool
	writer  io.Writer
}

// NewTimingLogger creates a new timing logger
func NewTimingLogger(prefix string, enabled bool) *TimingLogger {
	return &TimingLogger{
		prefix:  prefix,
		enabled: enabled,
		writer:  os.Stderr,
	}
}

// Time runs a function and logs its duration
func (t *TimingLogger) Time(operation string, fn func() error) error {
	if !t.enabled {
		return fn()
	}

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	if err != nil {
		fmt.Fprintf(t.writer, "%s%s... %dms (error: %v)\n", t.prefix, operation, duration.Milliseconds(), err)
	} else {
		fmt.Fprintf(t.writer, "%s%s... %dms\n", t.prefix, operation, duration.Milliseconds())
	}
	return err
}

// TimeValue runs a function and returns its value along with timing
func (t *TimingLogger) TimeValue(operation string, fn func() (bool, error)) (bool, error) {
	if !t.enabled {
		return fn()
	}

	start := time.Now()
	result, err := fn()
	duration := time.Since(start)

	if err != nil {
		fmt.Fprintf(t.writer, "%s%s... %dms (error: %v)\n", t.prefix, operation, duration.Milliseconds(), err)
	} else {
		fmt.Fprintf(t.writer, "%s%s... %dms\n", t.prefix, operation, duration.Milliseconds())
	}
	return result, err
}

// CheckStrategy defines the interface for checking update status
type CheckStrategy struct {
	git      GitRunner
	config   *JsonConfig
	syncDir  string
	logger    *TimingLogger
	timeouts  TimeoutConfig
}

// NewCheckStrategy creates a new check strategy with all optimizations
func NewCheckStrategy(git GitRunner, config *JsonConfig, syncDir string, logger *TimingLogger) *CheckStrategy {
	return &CheckStrategy{
		git:      git,
		config:   config,
		syncDir:  syncDir,
		logger:    logger,
		timeouts:  DefaultTimeouts(),
	}
}

func (c *CheckStrategy) CheckUnpushed() (bool, error) {
	var hasChanges bool
	err := c.logger.Time("Checking for unpushed changes", func() error {
		// Check git status for uncommitted changes (local, fast)
		ctx, cancel := context.WithTimeout(context.Background(), c.timeouts.GitLocal)
		defer cancel()

		cmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
		cmd.Dir = c.syncDir
		output, err := cmd.Output()
		if err != nil && ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("git status timeout")
		}
		if err != nil {
			return err
		}
		if len(output) > 0 {
			hasChanges = true
			return nil
		}

		// Check for unpushed commits (local, fast)
		cmd = exec.CommandContext(ctx, "git", "rev-list", "--left-right", "--count", "HEAD...@{u}")
		cmd.Dir = c.syncDir
		output, err = cmd.Output()
		if err != nil && ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("git rev-list timeout")
		}
		if err != nil {
			// Upstream not set, consider no unpushed
			return nil
		}

		parts := strings.Fields(string(output))
		if len(parts) >= 2 && parts[1] != "0" {
			hasChanges = true
		}
		return nil
	})
	return hasChanges, err
}

func (c *CheckStrategy) CheckUnpulled() (bool, error) {
	var hasChanges bool
	err := c.logger.Time("Checking for unpulled changes", func() error {
		// Get local HEAD (local, fast)
		ctx, cancel := context.WithTimeout(context.Background(), c.timeouts.GitLocal)
		defer cancel()

		cmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
		cmd.Dir = c.syncDir
		localHead, err := cmd.Output()
		if err != nil {
			// No commits yet
			return nil
		}
		localHash := strings.TrimSpace(string(localHead))

		// Get remote HEAD using ls-remote (remote, 5s timeout)
		ctxRemote, cancelRemote := context.WithTimeout(context.Background(), c.timeouts.GitRemote)
		defer cancelRemote()

		cmd = exec.CommandContext(ctxRemote, "git", "ls-remote", "--heads", "origin", "main")
		cmd.Dir = c.syncDir
		remoteHead, err := cmd.Output()
		if err != nil {
			if ctxRemote.Err() == context.DeadlineExceeded {
				return fmt.Errorf("git ls-remote timeout")
			}
			// Remote not reachable, can't check
			return nil
		}

		// ls-remote output format: "<hash>\trefs/heads/main"
		parts := strings.Split(string(remoteHead), "\t")
		if len(parts) < 2 {
			return nil
		}
		remoteHash := parts[0]

		// If hashes differ, there are unpulled changes
		if localHash != remoteHash {
			hasChanges = true
		}
		return nil
	})
	return hasChanges, err
}

func (c *CheckStrategy) CheckUnsyncedFiles() (bool, error) {
	var hasChanges bool
	err := c.logger.Time("Checking for unsynced source files", func() error {
		for tildePath := range c.config.Files {
			srcPath := ShorthandPath{}.New(tildePath)

			// Get source file hash
			srcHash, err := fileHash(srcPath.FullPath)
			if err != nil {
				// Source file doesn't exist or unreadable
				continue
			}

			// Get synced file hash
			hash := md5Hash(tildePath)
			baseName := filepath.Base(srcPath.FullPath)
			syncedPath := filepath.Join(c.syncDir, hash, baseName)

			syncedHash, err := fileHash(syncedPath)
			if err != nil {
				// Synced file doesn't exist - file is new
				hasChanges = true
			}

			// Compare hashes
			if srcHash != syncedHash {
				hasChanges = true
				return nil
			}
		}
		return nil
	})
	return hasChanges, err
}
