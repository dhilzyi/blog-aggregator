package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DB_url   string `json:"db_url"`
	Username string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (*Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return &Config{}, err
	}
	var data Config

	if configData, err := os.ReadFile(filePath); err == nil {
		if err := json.Unmarshal(configData, &data); err != nil {
			return &Config{}, err
		}
	} else {
		return &Config{}, err
	}

	return &data, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(homeDir, configFileName)

	return fullPath, nil
}

func Write(cfg *Config) error {
	dataFile, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, dataFile, 0644); err != nil {
		return err
	}

	return nil
}
