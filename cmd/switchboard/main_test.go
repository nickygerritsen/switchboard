package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nickygerritsen/switchboard/internal/logger"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `defaultBrowser: chrome
debug: false
logFile: ""
browsers:
  chrome:
    path: /usr/bin/chrome
rules:
  - match:
      - "*.google.com"
    browser: firefox
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Save original cfgFile and restore after test
	originalCfgFile := cfgFile
	defer func() { cfgFile = originalCfgFile }()

	// Test loading from specified file
	cfgFile = configPath
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}

	if cfg.DefaultBrowser != "chrome" {
		t.Errorf("loadConfig() DefaultBrowser = %q, want %q", cfg.DefaultBrowser, "chrome")
	}

	if len(cfg.Rules) != 1 {
		t.Errorf("loadConfig() rules count = %d, want 1", len(cfg.Rules))
	}
}

func TestLoadConfigInvalid(t *testing.T) {
	// Save original cfgFile and restore after test
	originalCfgFile := cfgFile
	defer func() { cfgFile = originalCfgFile }()

	// Test loading from non-existent file
	cfgFile = "/nonexistent/config.yaml"
	_, err := loadConfig()
	if err == nil {
		t.Error("loadConfig() should fail with non-existent file")
	}
}

func TestLoadConfigValidation(t *testing.T) {
	// Create a temporary invalid config file (missing defaultBrowser)
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	invalidConfig := `debug: false
logFile: ""
browsers:
  chrome:
    path: /usr/bin/chrome
`

	if err := os.WriteFile(configPath, []byte(invalidConfig), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Save original cfgFile and restore after test
	originalCfgFile := cfgFile
	defer func() { cfgFile = originalCfgFile }()

	cfgFile = configPath
	_, err := loadConfig()
	if err == nil {
		t.Error("loadConfig() should fail with invalid config")
	}
}

func TestCreateComponents(t *testing.T) {
	// Create a minimal valid config
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `defaultBrowser: chrome
debug: false
logFile: ""
browsers:
  chrome:
    path: /usr/bin/chrome
rules: []
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Save original cfgFile and restore after test
	originalCfgFile := cfgFile
	defer func() { cfgFile = originalCfgFile }()

	cfgFile = configPath
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}

	// Test creating components
	t.Run("createDetector", func(t *testing.T) {
		detector := createDetector(cfg)
		if detector == nil {
			t.Error("createDetector() returned nil")
		}
	})

	t.Run("createRouter", func(t *testing.T) {
		router := createRouter(cfg)
		if router == nil {
			t.Error("createRouter() returned nil")
		}
	})

	t.Run("createLauncher", func(t *testing.T) {
		launcher := createLauncher(cfg)
		if launcher == nil {
			t.Error("createLauncher() returned nil")
		}
	})
}

func TestInitLogger(t *testing.T) {
	// Create a minimal valid config
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// Create a minimal valid config
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `defaultBrowser: chrome
debug: true
logFile: ` + logFile + `
browsers:
  chrome:
    path: /usr/bin/chrome
rules: []
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Save original cfgFile and restore after test
	originalCfgFile := cfgFile
	defer func() { cfgFile = originalCfgFile }()

	cfgFile = configPath
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}

	// Test initializing logger
	err = initLogger(cfg)
	if err != nil {
		t.Fatalf("initLogger() failed: %v", err)
	}

	// Clean up
	defer func() { _ = logger.Close() }()

	// Verify log file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("initLogger() did not create log file")
	}
}
