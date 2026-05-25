package launcher

import (
	"fmt"
	"os/exec"

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

// Launch opens a URL in the specified browser with optional profile and incognito mode
func (l *Launcher) Launch(br *browser.Browser, url, profile string, incognito bool) error {
	if br == nil {
		return fmt.Errorf("browser cannot be nil")
	}
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	logger.Debug("Launching %s with URL: %s, profile: %s, incognito: %v", br.Name, url, profile, incognito)

	// Build command arguments
	args := buildArgs(br, url, profile, incognito)

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
func buildArgs(br *browser.Browser, url, profile string, incognito bool) []string {
	var args []string

	// Add incognito/private mode flag if requested
	if incognito {
		incognitoFlag := getIncognitoFlag(br.Name)
		if incognitoFlag != "" {
			args = append(args, incognitoFlag)
		} else if br.Name == "safari" {
			logger.Warn("Safari does not support private browsing via command-line, opening in normal mode")
		}
	}

	// Add profile-specific arguments if profile is specified and browser supports it
	if profile != "" && supportsProfiles(br.Name) {
		switch br.Name {
		case "chrome", "brave", "edge", "vivaldi", "chromium":
			// Chromium-based browsers use --profile-directory
			args = append(args, "--profile-directory="+profile)
		case "firefox":
			args = append(args, firefoxProfileArgs(br, profile)...)
		}
	}

	// Add the URL as the last argument
	args = append(args, url)

	return args
}

// firefoxProfileArgs returns the command-line arguments for selecting a
// Firefox profile by user-configured name. Old-style profiles (present in
// profiles.ini) are launched with `-P <IniName>`; new-style profiles
// (tracked only in the Profile Groups SQLite store) are launched with
// `--profile <Path>`, since `-P` cannot resolve them.
func firefoxProfileArgs(br *browser.Browser, profile string) []string {
	// Look for a matching detected profile, preferring an exact match on
	// the display Name (which may come from the SQLite store) and falling
	// back to the profiles.ini Name.
	var match *browser.Profile
	for i := range br.Profiles {
		if br.Profiles[i].Name == profile {
			match = &br.Profiles[i]
			break
		}
	}
	if match == nil {
		for i := range br.Profiles {
			if br.Profiles[i].IniName == profile {
				match = &br.Profiles[i]
				break
			}
		}
	}

	if match != nil {
		if match.IniName != "" {
			return []string{"-P", match.IniName}
		}
		if match.Path != "" {
			return []string{"--profile", match.Path}
		}
	}

	// No detected profile matched — keep the previous behaviour so users
	// without profile detection (or referring to a profile by its raw ini
	// name) continue to work.
	return []string{"-P", profile}
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

// getIncognitoFlag returns the incognito/private mode flag for a browser
func getIncognitoFlag(browserName string) string {
	switch browserName {
	case "chrome", "brave", "chromium", "vivaldi":
		return "--incognito"
	case "edge":
		return "--inprivate"
	case "firefox":
		return "--private-window"
	case "safari":
		// Safari doesn't support a command-line flag for private browsing
		// The best we can do is open normally
		return ""
	default:
		return ""
	}
}
