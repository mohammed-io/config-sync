package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type JsonConfig struct {
	Files       map[string]string `json:"files"`
	initialized bool
	folder      ShorthandPath
}

// checkInitialized returns an error if the config is not initialized
func (c *JsonConfig) checkInitialized() error {
	if !c.initialized {
		return errors.New("config not initialized. Call Initialize() first")
	}
	return nil
}

// Initialize loads or creates the config at the given folder
func (c *JsonConfig) Initialize(folder ShorthandPath) error {
	if err := os.MkdirAll(folder.Suffix("synced-files").FullPath, 0755); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	configPath := folder.Suffix("config.json")

	if _, err := os.Stat(filepath.Clean(configPath.FullPath)); errors.Is(err, os.ErrNotExist) {
		log.Printf("Config is not found at %s, initializing...\n", folder.TildePath)

		var formattedJson bytes.Buffer
		json.Indent(&formattedJson, []byte(`{"files": {}}`), "", "  ")

		if err := os.WriteFile(configPath.FullPath, formattedJson.Bytes(), 0755); err != nil {
			return err
		}
	} else {
		log.Printf("Config is detected at %s\n", folder.TildePath)
	}

	fileBytes, err := os.ReadFile(configPath.FullPath)
	if err != nil {
		return fmt.Errorf("could not read the file %s: %w", configPath, err)
	}

	if err := json.Unmarshal(fileBytes, c); err != nil {
		return fmt.Errorf("could not parse the json from the file %s: %w", configPath, err)
	}

	c.initialized = true
	c.folder = folder
	log.Printf("Config is initialized from %s\n", folder.TildePath)
	return nil
}

// Save writes the config to the config file
func (c *JsonConfig) Save() error {
	if err := c.checkInitialized(); err != nil {
		return err
	}
	jsonBytesToWrite, _ := json.MarshalIndent(c, "", "  ")
	return os.WriteFile(c.folder.Suffix("config.json").FullPath, jsonBytesToWrite, 0755)
}

// Track adds files to the config
func (c *JsonConfig) Track(files []string) error {
	if err := c.checkInitialized(); err != nil {
		return err
	}

	for _, file := range files {
		path := ShorthandPath{}.New(file)
		if _, err := os.Stat(path.FullPath); errors.Is(err, os.ErrNotExist) {
			log.Printf("Skipping %s: file does not exist\n", file)
			continue
		}

		if _, exists := c.Files[path.TildePath]; exists {
			log.Printf("Already tracked: %s\n", path.TildePath)
			continue
		}

		c.Files[path.TildePath] = filepath.Base(path.FullPath)
		log.Printf("Tracking: %s\n", path.TildePath)
	}

	return c.Save()
}

// Untrack removes files from the config
func (c *JsonConfig) Untrack(files []string) error {
	if err := c.checkInitialized(); err != nil {
		return err
	}

	for _, file := range files {
		path := ShorthandPath{}.New(file)
		if _, exists := c.Files[path.TildePath]; !exists {
			log.Printf("Not tracked: %s\n", path.TildePath)
			continue
		}

		delete(c.Files, path.TildePath)
		log.Printf("Untracked: %s\n", path.TildePath)
	}

	return c.Save()
}

// md5Hash returns the MD5 hash of a string
func md5Hash(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

// copyDir recursively copies a directory from src to dst
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// SyncFiles copies tracked files to the synced-files folder
func (c *JsonConfig) SyncFiles() error {
	if err := c.checkInitialized(); err != nil {
		return err
	}

	syncDir := c.folder.Suffix("synced-files")

	// Clean synced-files folder (remove all contents)
	if err := c.cleanSyncedFolder(syncDir.FullPath); err != nil {
		return fmt.Errorf("failed to clean synced folder: %w", err)
	}

	// Copy each tracked file to its MD5-hashed subfolder
	for tildePath := range c.Files {
		srcPath := ShorthandPath{}.New(tildePath)
		hash := md5Hash(tildePath)
		destDir := filepath.Join(syncDir.FullPath, hash)

		srcInfo, err := os.Stat(srcPath.FullPath)
		if err != nil {
			return fmt.Errorf("failed to stat %s: %w", tildePath, err)
		}

		if srcInfo.IsDir() {
			// For directories, copy the entire contents
			destPath := filepath.Join(destDir, filepath.Base(srcPath.FullPath))
			if err := copyDir(srcPath.FullPath, destPath); err != nil {
				return fmt.Errorf("failed to copy directory %s: %w", tildePath, err)
			}
			log.Printf("Synced directory: %s -> %s/%s\n", tildePath, hash, filepath.Base(srcPath.FullPath))
		} else {
			// For files, copy the file
			if err := os.MkdirAll(destDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destDir, err)
			}
			destPath := filepath.Join(destDir, filepath.Base(srcPath.FullPath))
			if err := copyFile(srcPath.FullPath, destPath); err != nil {
				return fmt.Errorf("failed to copy %s: %w", tildePath, err)
			}
			log.Printf("Synced file: %s -> %s/%s\n", tildePath, hash, filepath.Base(srcPath.FullPath))
		}
	}

	return nil
}

// cleanSyncedFolder removes all contents of the synced-files folder
func (c *JsonConfig) cleanSyncedFolder(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		if err := os.RemoveAll(fullPath); err != nil {
			return fmt.Errorf("failed to remove %s: %w", fullPath, err)
		}
	}

	return nil
}

// RestoreFiles copies tracked files from synced-files back to their original locations
func (c *JsonConfig) RestoreFiles() error {
	if err := c.checkInitialized(); err != nil {
		return err
	}

	syncDir := c.folder.Suffix("synced-files")

	for tildePath := range c.Files {
		destPath := ShorthandPath{}.New(tildePath)
		hash := md5Hash(tildePath)
		srcDir := filepath.Join(syncDir.FullPath, hash)

		// Check if hash folder exists
		if _, err := os.Stat(srcDir); err != nil {
			if os.IsNotExist(err) {
				log.Printf("Skipping %s: not found in synced-files\n", tildePath)
				continue
			}
			return fmt.Errorf("failed to stat sync dir for %s: %w", tildePath, err)
		}

		// Create parent directory if needed
		if err := os.MkdirAll(filepath.Dir(destPath.FullPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent dir for %s: %w", tildePath, err)
		}

		// Check what's inside the hash folder to determine if it's a file or directory
		baseName := filepath.Base(destPath.FullPath)
		srcPath := filepath.Join(srcDir, baseName)

		srcInfo, err := os.Stat(srcPath)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("Skipping %s: source not found\n", tildePath)
				continue
			}
			return fmt.Errorf("failed to stat source for %s: %w", tildePath, err)
		}

		if srcInfo.IsDir() {
			// For directories, copy the entire directory
			// Remove destination first if it exists as a file
			if destInfo, err := os.Stat(destPath.FullPath); err == nil && !destInfo.IsDir() {
				os.Remove(destPath.FullPath)
			}
			if err := copyDir(srcPath, destPath.FullPath); err != nil {
				return fmt.Errorf("failed to restore directory %s: %w", tildePath, err)
			}
			log.Printf("Restored directory: %s\n", tildePath)
		} else {
			// For files, copy the file
			// Remove destination first if it exists as a directory
			if destInfo, err := os.Stat(destPath.FullPath); err == nil && destInfo.IsDir() {
				os.RemoveAll(destPath.FullPath)
			}
			if err := copyFile(srcPath, destPath.FullPath); err != nil {
				return fmt.Errorf("failed to restore file %s: %w", tildePath, err)
			}
			log.Printf("Restored file: %s\n", tildePath)
		}
	}

	return nil
}
