package launcher

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nickygerritsen/switchboard/internal/browser"
	"github.com/nickygerritsen/switchboard/internal/config"
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

			args := buildArgs(br, tt.url, tt.profile)

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
	// Create a fake browser script for testing
	tmpDir := t.TempDir()
	browserPath := filepath.Join(tmpDir, "fake-browser")

	// Create a script that just echoes the arguments
	script := `#!/bin/sh
echo "$@" > ` + filepath.Join(tmpDir, "args.txt")

	if err := os.WriteFile(browserPath, []byte(script), 0755); err != nil {
		t.Fatalf("Failed to create fake browser: %v", err)
	}

	cfg := &config.Config{
		DefaultBrowser: "chrome",
	}

	launcher := NewLauncher(cfg)

	br := &browser.Browser{
		Name: "chrome",
		Path: browserPath,
	}

	err := launcher.Launch(br, "https://google.com", "")
	if err != nil {
		t.Fatalf("Launch() failed: %v", err)
	}

	// Note: We can't easily verify the browser was actually launched
	// without mocking exec.Command, but we can at least verify no errors
}

func TestLauncher_LaunchWithProfile(t *testing.T) {
	// Create a fake browser script
	tmpDir := t.TempDir()
	browserPath := filepath.Join(tmpDir, "fake-chrome")

	script := `#!/bin/sh
echo "$@" > ` + filepath.Join(tmpDir, "args.txt")

	if err := os.WriteFile(browserPath, []byte(script), 0755); err != nil {
		t.Fatalf("Failed to create fake browser: %v", err)
	}

	cfg := &config.Config{
		DefaultBrowser: "chrome",
	}

	launcher := NewLauncher(cfg)

	br := &browser.Browser{
		Name: "chrome",
		Path: browserPath,
	}

	err := launcher.Launch(br, "https://google.com", "Work")
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

	err := launcher.Launch(br, "https://google.com", "")
	if err == nil {
		t.Error("Launch() should fail with nonexistent browser path")
	}
}
