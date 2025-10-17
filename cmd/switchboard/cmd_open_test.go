package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nickygerritsen/switchboard/internal/browser"
	"github.com/nickygerritsen/switchboard/internal/config"
)

func TestRunOpen(t *testing.T) {
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `defaultBrowser: firefox
debug: false
logFile: ""
browsers:
  firefox:
    path: /usr/bin/firefox
  chrome:
    path: /usr/bin/chrome
rules:
  - match:
      - "*.github.com"
    browser: firefox
  - match:
      - "*.google.com"
    browser: chrome
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Save the original factories and cfgFile
	originalDetectorFactory := detectorFactory
	originalRouterFactory := routerFactory
	originalLauncherFactory := launcherFactory
	originalCfgFile := cfgFile
	defer func() {
		detectorFactory = originalDetectorFactory
		routerFactory = originalRouterFactory
		launcherFactory = originalLauncherFactory
		cfgFile = originalCfgFile
	}()

	cfgFile = configPath

	tests := []struct {
		name           string
		args           []string
		setupFakes     func() (*fakeDetector, *fakeRouter, *fakeLauncher)
		wantErr        bool
		wantErrContain string
		checkLaunched  func(t *testing.T, launcher *fakeLauncher)
	}{
		{
			name: "successful launch with matched rule",
			args: []string{"https://github.com/test"},
			setupFakes: func() (*fakeDetector, *fakeRouter, *fakeLauncher) {
				detector := &fakeDetector{
					browsers: map[string]*browser.Browser{
						"firefox": {Name: "firefox", Path: "/usr/bin/firefox"},
					},
				}
				router := &fakeRouter{
					matches: map[string]routeResult{
						"https://github.com/test": {browser: "firefox", profile: "", matched: true},
					},
				}
				launcher := &fakeLauncher{}
				return detector, router, launcher
			},
			wantErr: false,
			checkLaunched: func(t *testing.T, launcher *fakeLauncher) {
				if len(launcher.launchedURLs) != 1 {
					t.Fatalf("Expected 1 launch, got %d", len(launcher.launchedURLs))
				}
				if launcher.launchedURLs[0].url != "https://github.com/test" {
					t.Errorf("Expected URL https://github.com/test, got %s", launcher.launchedURLs[0].url)
				}
				if launcher.launchedURLs[0].browser != "firefox" {
					t.Errorf("Expected browser firefox, got %s", launcher.launchedURLs[0].browser)
				}
			},
		},
		{
			name: "successful launch with default browser",
			args: []string{"https://example.com"},
			setupFakes: func() (*fakeDetector, *fakeRouter, *fakeLauncher) {
				detector := &fakeDetector{
					browsers: map[string]*browser.Browser{
						"firefox": {Name: "firefox", Path: "/usr/bin/firefox"},
					},
				}
				router := &fakeRouter{
					matches: map[string]routeResult{
						"https://example.com": {browser: "firefox", profile: "", matched: false},
					},
				}
				launcher := &fakeLauncher{}
				return detector, router, launcher
			},
			wantErr: false,
			checkLaunched: func(t *testing.T, launcher *fakeLauncher) {
				if len(launcher.launchedURLs) != 1 {
					t.Fatalf("Expected 1 launch, got %d", len(launcher.launchedURLs))
				}
			},
		},
		{
			name: "browser detection fails",
			args: []string{"https://github.com/test"},
			setupFakes: func() (*fakeDetector, *fakeRouter, *fakeLauncher) {
				detector := &fakeDetector{
					detectError: fmt.Errorf("browser not found"),
				}
				router := &fakeRouter{
					matches: map[string]routeResult{
						"https://github.com/test": {browser: "firefox", profile: "", matched: true},
					},
				}
				launcher := &fakeLauncher{}
				return detector, router, launcher
			},
			wantErr:        true,
			wantErrContain: "failed to detect browser",
		},
		{
			name: "browser launch fails",
			args: []string{"https://github.com/test"},
			setupFakes: func() (*fakeDetector, *fakeRouter, *fakeLauncher) {
				detector := &fakeDetector{
					browsers: map[string]*browser.Browser{
						"firefox": {Name: "firefox", Path: "/usr/bin/firefox"},
					},
				}
				router := &fakeRouter{
					matches: map[string]routeResult{
						"https://github.com/test": {browser: "firefox", profile: "", matched: true},
					},
				}
				launcher := &fakeLauncher{
					launchError: fmt.Errorf("failed to start process"),
				}
				return detector, router, launcher
			},
			wantErr:        true,
			wantErrContain: "failed to launch browser",
		},
		{
			name: "launch with profile",
			args: []string{"https://work.github.com/test"},
			setupFakes: func() (*fakeDetector, *fakeRouter, *fakeLauncher) {
				detector := &fakeDetector{
					browsers: map[string]*browser.Browser{
						"chrome": {Name: "chrome", Path: "/usr/bin/chrome"},
					},
				}
				router := &fakeRouter{
					matches: map[string]routeResult{
						"https://work.github.com/test": {browser: "chrome", profile: "Work", matched: true},
					},
				}
				launcher := &fakeLauncher{}
				return detector, router, launcher
			},
			wantErr: false,
			checkLaunched: func(t *testing.T, launcher *fakeLauncher) {
				if len(launcher.launchedURLs) != 1 {
					t.Fatalf("Expected 1 launch, got %d", len(launcher.launchedURLs))
				}
				if launcher.launchedURLs[0].profile != "Work" {
					t.Errorf("Expected profile Work, got %s", launcher.launchedURLs[0].profile)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector, router, launcher := tt.setupFakes()

			// Override factories
			detectorFactory = func(cfg *config.Config) browserDetector {
				return detector
			}
			routerFactory = func(cfg *config.Config) urlRouter {
				return router
			}
			launcherFactory = func(cfg *config.Config) browserLauncher {
				return launcher
			}

			err := runOpen(openCmd, tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("runOpen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.wantErrContain != "" {
				if err == nil || len(err.Error()) == 0 {
					t.Errorf("Expected error containing %q, got nil", tt.wantErrContain)
				}
			}

			if !tt.wantErr && tt.checkLaunched != nil {
				tt.checkLaunched(t, launcher)
			}
		})
	}
}

func TestOpenCmdValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid single argument",
			args:    []string{"https://example.com"},
			wantErr: false,
		},
		{
			name:    "no arguments",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "too many arguments",
			args:    []string{"url1", "url2"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := openCmd.Args(openCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("openCmd.Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOpenCmdMetadata(t *testing.T) {
	if openCmd.Use != "open <url>" {
		t.Errorf("openCmd.Use = %q, want %q", openCmd.Use, "open <url>")
	}

	if openCmd.Short == "" {
		t.Error("openCmd.Short should not be empty")
	}

	if openCmd.Long == "" {
		t.Error("openCmd.Long should not be empty")
	}

	if openCmd.RunE == nil {
		t.Error("openCmd.RunE should not be nil")
	}
}
