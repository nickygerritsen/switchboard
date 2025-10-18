package launcher

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/nickygerritsen/switchboard/internal/browser"
	"github.com/nickygerritsen/switchboard/internal/config"
	"github.com/nickygerritsen/switchboard/internal/logger"
)

func TestBuildArgs(t *testing.T) {
	tests := []struct {
		name        string
		browserName string
		profile     string
		url         string
		wantProfile bool
	}{
		{
			name:        "chrome without profile",
			browserName: "chrome",
			profile:     "",
			url:         "https://google.com",
			wantProfile: false,
		},
		{
			name:        "chrome with profile",
			browserName: "chrome",
			profile:     "Work",
			url:         "https://google.com",
			wantProfile: true,
		},
		{
			name:        "firefox without profile",
			browserName: "firefox",
			profile:     "",
			url:         "https://google.com",
			wantProfile: false,
		},
		{
			name:        "firefox with profile",
			browserName: "firefox",
			profile:     "default",
			url:         "https://google.com",
			wantProfile: true,
		},
		{
			name:        "brave with profile",
			browserName: "brave",
			profile:     "Personal",
			url:         "https://google.com",
			wantProfile: true,
		},
		{
			name:        "vivaldi with profile",
			browserName: "vivaldi",
			profile:     "Work",
			url:         "https://google.com",
			wantProfile: true,
		},
		{
			name:        "safari ignores profile",
			browserName: "safari",
			profile:     "Work",
			url:         "https://google.com",
			wantProfile: false,
		},
		{
			name:        "edge with profile",
			browserName: "edge",
			profile:     "Work",
			url:         "https://google.com",
			wantProfile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := &browser.Browser{
				Name: tt.browserName,
				Path: "/usr/bin/" + tt.browserName,
			}

			args := buildArgs(br, tt.url, tt.profile, false)

			// Check URL is present
			found := false
			for _, arg := range args {
				if arg == tt.url {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("buildArgs() URL not found in args: %v", args)
			}

			// Check profile flag if expected
			if tt.wantProfile {
				foundProfile := false
				expectedChrome := "--profile-directory=" + tt.profile
				for i, arg := range args {
					// Chrome-based: --profile-directory=Profile
					if arg == expectedChrome {
						foundProfile = true
						break
					}
					// Firefox: -P Profile (two separate args)
					if arg == "-P" && i+1 < len(args) && args[i+1] == tt.profile {
						foundProfile = true
						break
					}
				}
				if !foundProfile {
					t.Errorf("buildArgs() profile flag not found in args: %v", args)
				}
			}
		})
	}
}

func TestSupportsProfiles(t *testing.T) {
	tests := []struct {
		name        string
		browserName string
		want        bool
	}{
		{"chrome supports profiles", "chrome", true},
		{"firefox supports profiles", "firefox", true},
		{"brave supports profiles", "brave", true},
		{"vivaldi supports profiles", "vivaldi", true},
		{"edge supports profiles", "edge", true},
		{"chromium supports profiles", "chromium", true},
		{"safari does not support profiles", "safari", false},
		{"unknown browser does not support profiles", "unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := supportsProfiles(tt.browserName)
			if got != tt.want {
				t.Errorf("supportsProfiles(%q) = %v, want %v", tt.browserName, got, tt.want)
			}
		})
	}
}

func TestNewLauncher(t *testing.T) {
	cfg := &config.Config{
		DefaultBrowser: "chrome",
	}

	launcher := NewLauncher(cfg)
	if launcher == nil {
		t.Fatal("NewLauncher() returned nil")
	}
	if launcher.config != cfg {
		t.Error("NewLauncher() did not store config")
	}
}

func TestLauncher_Launch(t *testing.T) {
	// Create a fake browser for testing
	tmpDir := t.TempDir()
	browserPath := createFakeBrowser(t, tmpDir, "fake-browser")

	cfg := &config.Config{
		DefaultBrowser: "chrome",
	}

	launcher := NewLauncher(cfg)

	br := &browser.Browser{
		Name: "chrome",
		Path: browserPath,
	}

	err := launcher.Launch(br, "https://google.com", "", false)
	if err != nil {
		t.Fatalf("Launch() failed: %v", err)
	}

	// Note: We can't easily verify the browser was actually launched
	// without mocking exec.Command, but we can at least verify no errors
}

func TestLauncher_LaunchWithProfile(t *testing.T) {
	// Create a fake browser for testing
	tmpDir := t.TempDir()
	browserPath := createFakeBrowser(t, tmpDir, "fake-chrome")

	cfg := &config.Config{
		DefaultBrowser: "chrome",
	}

	launcher := NewLauncher(cfg)

	br := &browser.Browser{
		Name: "chrome",
		Path: browserPath,
	}

	err := launcher.Launch(br, "https://google.com", "Work", false)
	if err != nil {
		t.Fatalf("Launch() with profile failed: %v", err)
	}
}

func TestLauncher_LaunchInvalidBrowser(t *testing.T) {
	cfg := &config.Config{
		DefaultBrowser: "chrome",
	}

	launcher := NewLauncher(cfg)

	br := &browser.Browser{
		Name: "chrome",
		Path: "/nonexistent/browser",
	}

	err := launcher.Launch(br, "https://google.com", "", false)
	if err == nil {
		t.Error("Launch() should fail with nonexistent browser path")
	}
}

func TestGetIncognitoFlag(t *testing.T) {
	tests := []struct {
		name        string
		browserName string
		want        string
	}{
		{"chrome has --incognito", "chrome", "--incognito"},
		{"brave has --incognito", "brave", "--incognito"},
		{"chromium has --incognito", "chromium", "--incognito"},
		{"vivaldi has --incognito", "vivaldi", "--incognito"},
		{"edge has --inprivate", "edge", "--inprivate"},
		{"firefox has --private-window", "firefox", "--private-window"},
		{"safari has no flag", "safari", ""},
		{"unknown browser has no flag", "unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getIncognitoFlag(tt.browserName)
			if got != tt.want {
				t.Errorf("getIncognitoFlag(%q) = %q, want %q", tt.browserName, got, tt.want)
			}
		})
	}
}

func TestBuildArgsWithIncognito(t *testing.T) {
	tests := []struct {
		name         string
		browserName  string
		incognito    bool
		wantFlag     string
		wantHasFlag  bool
	}{
		{
			name:         "chrome with incognito",
			browserName:  "chrome",
			incognito:    true,
			wantFlag:     "--incognito",
			wantHasFlag:  true,
		},
		{
			name:         "chrome without incognito",
			browserName:  "chrome",
			incognito:    false,
			wantFlag:     "--incognito",
			wantHasFlag:  false,
		},
		{
			name:         "firefox with incognito",
			browserName:  "firefox",
			incognito:    true,
			wantFlag:     "--private-window",
			wantHasFlag:  true,
		},
		{
			name:         "edge with incognito",
			browserName:  "edge",
			incognito:    true,
			wantFlag:     "--inprivate",
			wantHasFlag:  true,
		},
		{
			name:         "safari with incognito request",
			browserName:  "safari",
			incognito:    true,
			wantFlag:     "",
			wantHasFlag:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := &browser.Browser{
				Name: tt.browserName,
				Path: "/usr/bin/" + tt.browserName,
			}

			args := buildArgs(br, "https://google.com", "", tt.incognito)

			// Check if incognito flag is present
			found := false
			for _, arg := range args {
				if arg == tt.wantFlag {
					found = true
					break
				}
			}

			if found != tt.wantHasFlag {
				if tt.wantHasFlag {
					t.Errorf("buildArgs() should include %q flag but didn't, args: %v", tt.wantFlag, args)
				} else {
					t.Errorf("buildArgs() should not include incognito flag but did, args: %v", args)
				}
			}
		})
	}
}

func TestSafariIncognitoWarning(t *testing.T) {
	// Setup logger with temp file to capture log output
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	cfg := &config.Config{
		DefaultBrowser: "chrome",
		Debug:          false, // Warn level to capture WARN messages
		LogFile:        logFile,
	}

	if err := logger.Init(cfg); err != nil {
		t.Fatalf("Failed to init logger: %v", err)
	}
	defer func() { _ = logger.Close() }()

	br := &browser.Browser{
		Name: "safari",
		Path: "/Applications/Safari.app/Contents/MacOS/Safari",
	}

	url := "https://example.com"
	args := buildArgs(br, url, "", true) // incognito=true

	// Close logger to flush to file
	if err := logger.Close(); err != nil {
		t.Fatalf("Failed to close logger: %v", err)
	}

	// Read log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify the warning was logged
	expectedWarning := "Safari does not support private browsing via command-line, opening in normal mode"
	if !strings.Contains(logContent, expectedWarning) {
		t.Errorf("Expected warning message not found in logs.\nExpected: %s\nLog content: %s", expectedWarning, logContent)
	}

	// Verify no incognito-like flags are present
	for _, arg := range args {
		if arg == "--incognito" || arg == "--private-window" || arg == "--inprivate" {
			t.Errorf("buildArgs() should not include incognito flags for Safari, but found: %s", arg)
		}
	}

	// Verify URL is still present
	found := false
	for _, arg := range args {
		if arg == url {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("buildArgs() should include URL %s in args: %v", url, args)
	}
}

// createFakeBrowser creates a platform-specific fake browser executable for testing
func createFakeBrowser(t *testing.T, dir, name string) string {
	t.Helper()

	var browserPath string
	var script string

	if runtime.GOOS == "windows" {
		// Windows: create a batch file
		browserPath = filepath.Join(dir, name+".bat")
		script = "@echo off\nrem Fake browser for testing\nexit /b 0\n"
	} else {
		// Unix (Linux/macOS): create a shell script
		browserPath = filepath.Join(dir, name)
		script = "#!/bin/sh\n# Fake browser for testing\nexit 0\n"
	}

	if err := os.WriteFile(browserPath, []byte(script), 0755); err != nil {
		t.Fatalf("Failed to create fake browser: %v", err)
	}

	return browserPath
}
