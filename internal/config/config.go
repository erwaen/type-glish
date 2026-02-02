package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = "type-glish-config.json"

type Config struct {
	Provider     string `json:"provider"`
	GeminiAPIKey string `json:"gemini_api_key"`
	GeminiModel  string `json:"gemini_model"`
}

func GetConfigPath() (string, error) {
	// Try to use UserConfigDir first
	configDir, err := os.UserConfigDir()
	if err == nil {
		appDir := filepath.Join(configDir, "type-glish")
		// Ensure directory exists
		if _, err := os.Stat(appDir); os.IsNotExist(err) {
			os.MkdirAll(appDir, 0755)
		}
		return filepath.Join(appDir, configFileName), nil
	}

	// Fallback to current directory
	return configFileName, nil
}

func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Returns default config if not found
			return &Config{
				Provider: "", // Default to empty to trigger Menu
			}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
