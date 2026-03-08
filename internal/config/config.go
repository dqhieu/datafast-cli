package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey    string `json:"api_key"`
	WebsiteID string `json:"website_id"`
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "datafast")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

func Load() (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(configPath())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	if val := os.Getenv("DATAFAST_API_KEY"); val != "" {
		cfg.APIKey = val
	}
	if val := os.Getenv("DATAFAST_WEBSITE_ID"); val != "" {
		cfg.WebsiteID = val
	}

	return cfg, nil
}

func GetAPIKey() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}
	if cfg.APIKey == "" {
		return "", errors.New("API key not configured. Run 'datafast auth login' or set DATAFAST_API_KEY")
	}
	return cfg.APIKey, nil
}

func GetWebsiteID() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}
	if cfg.WebsiteID == "" {
		return "", errors.New("website ID not configured. Run 'datafast auth login' or set DATAFAST_WEBSITE_ID")
	}
	return cfg.WebsiteID, nil
}

func Save(cfg *Config) error {
	if err := os.MkdirAll(configDir(), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0600)
}

func Clear() error {
	err := os.Remove(configPath())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
