package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// configPathFunc is a function type for getting the config path
type configPathFunc func() (string, error)

// configPathProvider is the function used to get the config path
// Can be overridden in tests for dependency injection
var configPathProvider configPathFunc = getDefaultConfigPath

// GetConfigPath returns the configuration file path
// Uses configPathProvider which can be overridden in tests
func GetConfigPath() (string, error) {
	return configPathProvider()
}

// SetConfigPathProvider allows overriding the config path function for testing
// Returns the previous provider so it can be restored
func SetConfigPathProvider(provider configPathFunc) configPathFunc {
	old := configPathProvider
	configPathProvider = provider
	return old
}

// getDefaultConfigPath returns the default configuration file path for the current OS
func getDefaultConfigPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

// GetConfigDir returns the configuration directory for the current OS
func getConfigDir() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "darwin", "linux":
		// Use XDG_CONFIG_HOME if set, otherwise use ~/.config
		if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
			configDir = filepath.Join(xdgConfig, "switchboard")
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			configDir = filepath.Join(home, ".config", "switchboard")
		}
	case "windows":
		// Use APPDATA on Windows
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		configDir = filepath.Join(appData, "switchboard")
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return configDir, nil
}

// ensureConfigDir creates the configuration directory if it doesn't exist
func ensureConfigDir() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return nil
}

// GetLogPath returns the default log file path for the current OS
func GetLogPath() (string, error) {
	var logDir string

	switch runtime.GOOS {
	case "darwin":
		// macOS: ~/Library/Logs/switchboard
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		logDir = filepath.Join(home, "Library", "Logs", "switchboard")
	case "linux":
		// Linux: ~/.local/state/switchboard or ~/.cache/switchboard
		if xdgState := os.Getenv("XDG_STATE_HOME"); xdgState != "" {
			logDir = filepath.Join(xdgState, "switchboard")
		} else if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
			logDir = filepath.Join(xdgCache, "switchboard")
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			logDir = filepath.Join(home, ".local", "state", "switchboard")
		}
	case "windows":
		// Windows: %LOCALAPPDATA%\switchboard\logs
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return "", fmt.Errorf("LOCALAPPDATA environment variable not set")
		}
		logDir = filepath.Join(localAppData, "switchboard", "logs")
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Ensure log directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create log directory: %w", err)
	}

	return filepath.Join(logDir, "switchboard.log"), nil
}
