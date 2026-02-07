package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type JsonConfig struct {
	Files map[string]string `json:"files"`
}

// configFolder returns the ShorthandPath for ~/.config-sync
func configFolder() ShorthandPath {
	return ShorthandPath{}.New("~/.config-sync")
}

// syncedFilesFolder returns the ShorthandPath for the synced-files directory
func syncedFilesFolder() ShorthandPath {
	return configFolder().Suffix("synced-files")
}

// configPath returns the ShorthandPath for the config.json file
func configPath() ShorthandPath {
	return configFolder().Suffix("config.json")
}

func loadConfigFromJson(configPath ShorthandPath, jsonConfig *JsonConfig) error {
	fileBytes, err := os.ReadFile(configPath.FullPath)

	if err != nil {
		return fmt.Errorf("Could not read the file %s", configPath)
	}

	err = json.Unmarshal(fileBytes, &jsonConfig)

	if err != nil {
		return fmt.Errorf("Could not parse the json from the file %s", configPath)
	}

	return nil
}

func initializeConfig() (JsonConfig, error) {
	var jsonConfig JsonConfig

	err := os.MkdirAll(syncedFilesFolder().FullPath, 0755)

	if err != nil && !errors.Is(err, os.ErrExist) {
		return jsonConfig, err
	}

	configPath := configPath()

	if _, err := os.Stat(filepath.Clean(configPath.FullPath)); errors.Is(err, os.ErrNotExist) {
		log.Printf("Config is not found at %s, initializing...\n", configFolder().TildePath)

		var formattedJson bytes.Buffer
		json.Indent(&formattedJson, []byte(`{"files": {}}`), "", "  ")

		os.WriteFile(configPath.FullPath, formattedJson.Bytes(), 0755)
	} else {
		log.Printf("Config is detected at %s\n", configFolder().TildePath)
	}

	err = loadConfigFromJson(configPath, &jsonConfig)

	if err != nil {
		return jsonConfig, err
	}

	log.Printf("Config is initialized from %s\n", configFolder().TildePath)

	return jsonConfig, nil
}

func trackFileInConfig(config *JsonConfig, pathToFile ShorthandPath) error {
	config.Files[pathToFile.TildePath] = filepath.Base(pathToFile.FullPath)

	jsonBytesToWrite, _ := json.MarshalIndent(config, "", "  ")

	err := os.WriteFile(configPath().FullPath, jsonBytesToWrite, 0755)

	return err
}

func untrackFileInConfig(config *JsonConfig, pathToFile ShorthandPath) error {
	delete(config.Files, pathToFile.TildePath)

	jsonBytesToWrite, _ := json.MarshalIndent(config, "", "  ")

	err := os.WriteFile(configPath().FullPath, jsonBytesToWrite, 0755)

	return err
}

func main() {
	jsonConfig, err := initializeConfig()

	if err != nil {
		log.Fatalf("Could not inialize the config: %s", err)
	}

	fmt.Printf("%s", jsonConfig)
}
