package launcher

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/nickygerritsen/switchboard/internal/browser"
	"github.com/nickygerritsen/switchboard/internal/config"
	"github.com/nickygerritsen/switchboard/internal/logger"
)

// Launcher handles launching browsers with URLs
type Launcher struct {
	config *config.Config
}

// NewLauncher creates a new launcher instance
func NewLauncher(cfg *config.Config) *Launcher {
	return &Launcher{
		config: cfg,
	}
}

// Launch opens a URL in the specified browser with optional profile
func (l *Launcher) Launch(br *browser.Browser, url, profile string) error {
	if br == nil {
		return fmt.Errorf("browser cannot be nil")
	}
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	logger.Debug("Launching %s with URL: %s, profile: %s", br.Name, url, profile)

	// Build command arguments
	args := buildArgs(br, url, profile)

	// Create command
	cmd := exec.Command(br.Path, args...)

	// Start the browser in the background
	err := cmd.Start()
	if err != nil {
		logger.Error("Failed to launch %s: %v", br.Name, err)
		return fmt.Errorf("failed to launch browser: %w", err)
	}

	logger.Info("Successfully launched %s with URL: %s", br.Name, url)

	// Don't wait for the browser to exit - let it run in the background
	go func() {
		_ = cmd.Wait()
	}()

	return nil
}

// buildArgs constructs the command-line arguments for launching a browser
func buildArgs(br *browser.Browser, url, profile string) []string {
	var args []string

	// Add profile-specific arguments if profile is specified and browser supports it
	if profile != "" && supportsProfiles(br.Name) {
		switch br.Name {
		case "chrome", "brave", "edge", "vivaldi", "chromium":
			// Chromium-based browsers use --profile-directory
			args = append(args, "--profile-directory="+profile)
		case "firefox":
			// Firefox uses -P flag
			args = append(args, "-P", profile)
		}
	}

	// Add platform-specific flags
	args = append(args, getPlatformArgs(br.Name)...)

	// Add the URL as the last argument
	args = append(args, url)

	return args
}

// supportsProfiles returns true if the browser supports profiles
func supportsProfiles(browserName string) bool {
	switch browserName {
	case "chrome", "firefox", "brave", "vivaldi", "edge", "chromium":
		return true
	default:
		return false
	}
}

// getPlatformArgs returns platform-specific arguments for the browser
func getPlatformArgs(browserName string) []string {
	var args []string

	switch runtime.GOOS {
	case "darwin":
		// macOS-specific flags
		// For most browsers on macOS, we want to open a new window/tab
		// Safari is typically launched via 'open -a Safari' but we're using the direct binary path
		switch browserName {
		case "chrome", "brave", "edge", "vivaldi", "chromium":
			// No special flags needed - Chrome-based browsers handle this well
		case "firefox":
			// Firefox on macOS
			args = append(args, "--new-window")
		}

	case "linux":
		// Linux-specific flags
		switch browserName {
		case "chrome", "brave", "edge", "vivaldi", "chromium":
			// Chromium-based browsers
			args = append(args, "--new-window")
		case "firefox":
			// Firefox on Linux
			args = append(args, "--new-window")
		}

	case "windows":
		// Windows-specific flags
		switch browserName {
		case "chrome", "brave", "edge", "vivaldi", "chromium":
			// Chromium-based browsers
			args = append(args, "--new-window")
		case "firefox":
			// Firefox on Windows
			args = append(args, "-new-window")
		}
	}

	return args
}
