package registration

import (
	"fmt"
	"os"
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
// Note: This does NOT resolve symlinks to ensure stability with package managers
// like Homebrew that use symlinks to stable paths (e.g., /opt/homebrew/bin/switchboard)
// rather than versioned paths (e.g., /opt/homebrew/Cellar/switchboard/1.0.1/bin/switchboard)
func getBinaryPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	return exe, nil
}
