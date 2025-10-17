package registration

import (
	"fmt"
	"os"
	"path/filepath"
)

// Registrar handles browser registration for the current platform
type Registrar interface {
	// Register registers Switchboard as a browser/URL handler
	Register() error

	// Unregister removes Switchboard's browser registration
	Unregister() error

	// IsRegistered checks if Switchboard is currently registered
	IsRegistered() (bool, error)

	// GetBinaryPath returns the path to the switchboard binary
	GetBinaryPath() (string, error)
}

// getBinaryPath returns the absolute path to the current executable
func getBinaryPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve symlinks
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	return exe, nil
}
