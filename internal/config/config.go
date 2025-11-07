package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbURL string `json:"db_url"`
}

func Read() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}
	}
	cfgPath := homeDir + "/.gatorconfig.json"

	cfgContent, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}
	}

	var cfg Config
	if err := json.Unmarshal(cfgContent, &cfg); err != nil {
		return Config{}
	}

	fmt.Println(cfg)
	return cfg
}
