package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL string `json:"db_url"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}

func Read() (Config, error) {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	cfgContent, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(cfgContent, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func write(cfg Config) error {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(cfg, "", "	")
	if err != nil {
		return err
	}

	err = os.WriteFile(cfgPath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
