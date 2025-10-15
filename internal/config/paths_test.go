package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}

	if configPath == "" {
		t.Error("GetConfigPath() returned empty path")
	}

	if !strings.HasSuffix(configPath, "config.yaml") {
		t.Errorf("GetConfigPath() = %v, should end with 'config.yaml'", configPath)
	}

	// Check OS-specific expectations
	switch runtime.GOOS {
	case "darwin", "linux":
		if !strings.Contains(configPath, "switchboard") {
			t.Errorf("GetConfigPath() = %v, should contain 'switchboard'", configPath)
		}
	case "windows":
		if !strings.Contains(configPath, "switchboard") {
			t.Errorf("GetConfigPath() = %v, should contain 'switchboard'", configPath)
		}
	}
}

func TestGetConfigDir(t *testing.T) {
	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() failed: %v", err)
	}

	if configDir == "" {
		t.Error("getConfigDir() returned empty path")
	}

	if !strings.Contains(configDir, "switchboard") {
		t.Errorf("getConfigDir() = %v, should contain 'switchboard'", configDir)
	}
}

func TestGetConfigDirWithXDGConfigHome(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping XDG test on Windows")
	}

	// Save original environment
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if originalXDG != "" {
			os.Setenv("XDG_CONFIG_HOME", originalXDG)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
	}()

	// Set custom XDG_CONFIG_HOME
	testDir := "/tmp/test-config"
	os.Setenv("XDG_CONFIG_HOME", testDir)

	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() failed: %v", err)
	}

	expected := filepath.Join(testDir, "switchboard")
	if configDir != expected {
		t.Errorf("getConfigDir() = %v, want %v", configDir, expected)
	}
}

func TestGetLogPath(t *testing.T) {
	logPath, err := GetLogPath()
	if err != nil {
		t.Fatalf("GetLogPath() failed: %v", err)
	}

	if logPath == "" {
		t.Error("GetLogPath() returned empty path")
	}

	if !strings.HasSuffix(logPath, "switchboard.log") {
		t.Errorf("GetLogPath() = %v, should end with 'switchboard.log'", logPath)
	}

	// Check OS-specific expectations
	switch runtime.GOOS {
	case "darwin":
		if !strings.Contains(logPath, "Library/Logs") {
			t.Errorf("GetLogPath() on macOS = %v, should contain 'Library/Logs'", logPath)
		}
	case "linux":
		// Should contain either .local/state or .cache
		if !strings.Contains(logPath, ".local") && !strings.Contains(logPath, ".cache") {
			t.Errorf("GetLogPath() on Linux = %v, should contain '.local' or '.cache'", logPath)
		}
	case "windows":
		if !strings.Contains(logPath, "switchboard") {
			t.Errorf("GetLogPath() on Windows = %v, should contain 'switchboard'", logPath)
		}
	}

	// Verify the directory was created
	logDir := filepath.Dir(logPath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		t.Errorf("GetLogPath() should create log directory, but %v does not exist", logDir)
	}
}

func TestEnsureConfigDir(t *testing.T) {
	// This is tested implicitly through other tests, but let's test it explicitly
	err := ensureConfigDir()
	if err != nil {
		t.Fatalf("ensureConfigDir() failed: %v", err)
	}

	// Verify the directory exists
	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() failed: %v", err)
	}

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("ensureConfigDir() should create directory, but %v does not exist", configDir)
	}
}

func TestGetConfigDirUnsupportedOS(t *testing.T) {
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" && runtime.GOOS != "windows" {
		_, err := getConfigDir()
		if err == nil {
			t.Error("getConfigDir() should return error for unsupported OS")
		}
		if !strings.Contains(err.Error(), "unsupported operating system") {
			t.Errorf("getConfigDir() error = %v, should mention unsupported OS", err)
		}
	} else {
		t.Skip("Skipping unsupported OS test on supported OS")
	}
}
