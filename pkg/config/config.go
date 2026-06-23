package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Cookie         string `yaml:"cookie"`
	UserAgent      string `yaml:"user_agent"`
	OutputDir      string `yaml:"output_dir"`
	NamingTemplate string `yaml:"naming_template"`
	Count          int    `yaml:"count"`
	Proxy          string `yaml:"proxy"`
	Concurrency    int    `yaml:"concurrency"`
}

const (
	DefaultConfigPath = "config.yaml"
	defaultUserAgent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36"
)

// DefaultConfig returns a config struct with default values
func DefaultConfig() *Config {
	return &Config{
		Cookie:         "",
		UserAgent:      defaultUserAgent,
		OutputDir:      "./downloads",
		NamingTemplate: "{nickname}/{publish_date}_{title}",
		Count:          50,
		Proxy:          "",
		Concurrency:    5,
	}
}

// LoadConfig loads config from config.yaml, creating a default one if it doesn't exist
func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	// If config file doesn't exist, create default one
	if _, err := os.Stat(DefaultConfigPath); os.IsNotExist(err) {
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return cfg, err
		}
		_ = os.WriteFile(DefaultConfigPath, data, 0644)
		return cfg, nil
	}

	data, err := os.ReadFile(DefaultConfigPath)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(data, cfg)
	return cfg, err
}
