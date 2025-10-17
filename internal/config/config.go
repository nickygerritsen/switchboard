package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	DefaultBrowser string             `yaml:"defaultBrowser"`
	Debug          bool               `yaml:"debug"`
	LogFile        string             `yaml:"logFile"`
	Browsers       map[string]Browser `yaml:"browsers"`
	Rules          []Rule             `yaml:"rules"`
}

// Browser represents a browser configuration
type Browser struct {
	Path string `yaml:"path"`
}

// Rule represents a URL routing rule
type Rule struct {
	Match     []string `yaml:"match"`
	Browser   string   `yaml:"browser"`
	Profile   string   `yaml:"profile,omitempty"`
	Incognito bool     `yaml:"incognito,omitempty"`
	Rewrite   string   `yaml:"rewrite,omitempty"`
}

// Load loads the configuration from the default config file location
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	return LoadFrom(configPath)
}

// LoadFrom loads the configuration from the specified file
func LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.DefaultBrowser == "" {
		return fmt.Errorf("defaultBrowser is required")
	}

	// Check that all rules reference valid browsers
	for i, rule := range c.Rules {
		if rule.Browser == "" {
			return fmt.Errorf("rule %d: browser is required", i)
		}
		if len(rule.Match) == 0 {
			return fmt.Errorf("rule %d: at least one match pattern is required", i)
		}
	}

	return nil
}

// Save writes the configuration to the default config file location
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	return c.SaveTo(configPath)
}

// SaveTo writes the configuration to the specified file
func (c *Config) SaveTo(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure the directory exists
	if err := ensureConfigDir(); err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetDefaultConfig returns a default configuration
func GetDefaultConfig() *Config {
	return &Config{
		DefaultBrowser: "chrome",
		Debug:          false,
		LogFile:        "auto",
		Browsers:       make(map[string]Browser),
		Rules:          []Rule{},
	}
}
