package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				DefaultBrowser: "chrome",
				Rules: []Rule{
					{
						Match:   []string{"*.google.com"},
						Browser: "firefox",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing default browser",
			config: &Config{
				DefaultBrowser: "",
				Rules:          []Rule{},
			},
			wantErr: true,
		},
		{
			name: "rule missing browser",
			config: &Config{
				DefaultBrowser: "chrome",
				Rules: []Rule{
					{
						Match:   []string{"*.google.com"},
						Browser: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rule missing match patterns",
			config: &Config{
				DefaultBrowser: "chrome",
				Rules: []Rule{
					{
						Match:   []string{},
						Browser: "firefox",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid config with profile",
			config: &Config{
				DefaultBrowser: "chrome",
				Rules: []Rule{
					{
						Match:   []string{"*.company.com"},
						Browser: "chrome",
						Profile: "Work",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadFrom(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{
			name: "valid yaml",
			yaml: `defaultBrowser: chrome
debug: true
logFile: auto
rules:
  - match:
      - "*.google.com"
    browser: firefox
`,
			wantErr: false,
		},
		{
			name:    "invalid yaml",
			yaml:    `this is not valid yaml: [[[`,
			wantErr: true,
		},
		{
			name: "missing required field",
			yaml: `debug: true
rules: []
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "config.yaml")

			if err := os.WriteFile(tmpFile, []byte(tt.yaml), 0644); err != nil {
				t.Fatalf("Failed to create temp config file: %v", err)
			}

			_, err := LoadFrom(tmpFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFrom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveToAndLoadFrom(t *testing.T) {
	config := &Config{
		DefaultBrowser: "chrome",
		Debug:          true,
		LogFile:        "auto",
		Browsers: map[string]Browser{
			"firefox": {Path: "/usr/bin/firefox"},
		},
		Rules: []Rule{
			{
				Match:   []string{"*.google.com", "google.com"},
				Browser: "firefox",
			},
			{
				Match:   []string{"*.github.com"},
				Browser: "chrome",
				Profile: "Work",
			},
		},
	}

	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	// Save config
	if err := config.SaveTo(tmpFile); err != nil {
		t.Fatalf("SaveTo() failed: %v", err)
	}

	// Load config
	loaded, err := LoadFrom(tmpFile)
	if err != nil {
		t.Fatalf("LoadFrom() failed: %v", err)
	}

	// Compare
	if loaded.DefaultBrowser != config.DefaultBrowser {
		t.Errorf("DefaultBrowser = %v, want %v", loaded.DefaultBrowser, config.DefaultBrowser)
	}
	if loaded.Debug != config.Debug {
		t.Errorf("Debug = %v, want %v", loaded.Debug, config.Debug)
	}
	if len(loaded.Rules) != len(config.Rules) {
		t.Errorf("Rules length = %v, want %v", len(loaded.Rules), len(config.Rules))
	}
}

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()

	if config.DefaultBrowser == "" {
		t.Error("Default config should have a default browser")
	}
	if config.Browsers == nil {
		t.Error("Default config should have initialized Browsers map")
	}
	if config.Rules == nil {
		t.Error("Default config should have initialized Rules slice")
	}
}
