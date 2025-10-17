package router

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nickygerritsen/switchboard/internal/config"
)

func TestNewRouter(t *testing.T) {
	cfg := &config.Config{
		DefaultBrowser: "chrome",
	}

	router := NewRouter(cfg)
	if router == nil {
		t.Fatal("NewRouter() returned nil")
	}
	if router.config != cfg {
		t.Error("NewRouter() did not store config")
	}
}

func TestRouter_FindMatch(t *testing.T) {
	cfg := &config.Config{
		DefaultBrowser: "chrome",
		Rules: []config.Rule{
			{
				Match:   []string{"*.google.com"},
				Browser: "firefox",
			},
			{
				Match:   []string{"github.com"},
				Browser: "chrome",
				Profile: "Work",
			},
			{
				Match:   []string{"localhost:3000"},
				Browser: "brave",
			},
		},
	}

	router := NewRouter(cfg)

	tests := []struct {
		name        string
		url         string
		wantBrowser string
		wantProfile string
		wantMatched bool
	}{
		{
			name:        "matches google subdomain",
			url:         "https://mail.google.com",
			wantBrowser: "firefox",
			wantProfile: "",
			wantMatched: true,
		},
		{
			name:        "matches github with profile",
			url:         "https://github.com/user/repo",
			wantBrowser: "chrome",
			wantProfile: "Work",
			wantMatched: true,
		},
		{
			name:        "matches localhost with port",
			url:         "http://localhost:3000/app",
			wantBrowser: "brave",
			wantProfile: "",
			wantMatched: true,
		},
		{
			name:        "no match returns default",
			url:         "https://example.com",
			wantBrowser: "chrome",
			wantProfile: "",
			wantMatched: false,
		},
		{
			name:        "first match wins",
			url:         "https://mail.google.com",
			wantBrowser: "firefox",
			wantProfile: "",
			wantMatched: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			browser, profile, matched := router.FindMatch(tt.url)
			if browser != tt.wantBrowser {
				t.Errorf("FindMatch() browser = %q, want %q", browser, tt.wantBrowser)
			}
			if profile != tt.wantProfile {
				t.Errorf("FindMatch() profile = %q, want %q", profile, tt.wantProfile)
			}
			if matched != tt.wantMatched {
				t.Errorf("FindMatch() matched = %v, want %v", matched, tt.wantMatched)
			}
		})
	}
}

func TestRouter_FindMatchInvalidURL(t *testing.T) {
	cfg := &config.Config{
		DefaultBrowser: "chrome",
		Rules: []config.Rule{
			{
				Match:   []string{"google.com"},
				Browser: "firefox",
			},
		},
	}

	router := NewRouter(cfg)

	// Invalid URL should return default browser
	browser, profile, matched := router.FindMatch("not a url")
	if browser != "chrome" {
		t.Errorf("FindMatch() with invalid URL browser = %q, want %q", browser, "chrome")
	}
	if profile != "" {
		t.Errorf("FindMatch() with invalid URL profile = %q, want empty", profile)
	}
	if matched {
		t.Error("FindMatch() with invalid URL should not match")
	}
}

func TestRouter_RouteIntegration(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `defaultBrowser: chrome
debug: false
logFile: ""
browsers:
  chrome:
    path: /usr/bin/chrome
  firefox:
    path: /usr/bin/firefox
rules:
  - match:
      - "*.google.com"
    browser: firefox
  - match:
      - "github.com"
    browser: chrome
    profile: Work
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Load config
	cfg, err := config.LoadFrom(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	router := NewRouter(cfg)

	tests := []struct {
		name        string
		url         string
		wantBrowser string
		wantProfile string
	}{
		{
			name:        "route google to firefox",
			url:         "https://mail.google.com",
			wantBrowser: "firefox",
			wantProfile: "",
		},
		{
			name:        "route github to chrome with profile",
			url:         "https://github.com/user/repo",
			wantBrowser: "chrome",
			wantProfile: "Work",
		},
		{
			name:        "route unknown to default",
			url:         "https://example.com",
			wantBrowser: "chrome",
			wantProfile: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			browser, profile, _ := router.FindMatch(tt.url)
			if browser != tt.wantBrowser {
				t.Errorf("Route() browser = %q, want %q", browser, tt.wantBrowser)
			}
			if profile != tt.wantProfile {
				t.Errorf("Route() profile = %q, want %q", profile, tt.wantProfile)
			}
		})
	}
}

func TestRouter_FindMatchEmptyRules(t *testing.T) {
	cfg := &config.Config{
		DefaultBrowser: "chrome",
		Rules:          []config.Rule{},
	}

	router := NewRouter(cfg)

	browser, profile, matched := router.FindMatch("https://google.com")
	if browser != "chrome" {
		t.Errorf("FindMatch() with no rules browser = %q, want %q", browser, "chrome")
	}
	if profile != "" {
		t.Errorf("FindMatch() with no rules profile = %q, want empty", profile)
	}
	if matched {
		t.Error("FindMatch() with no rules should not match")
	}
}

func TestRouter_FindMatchMultiplePatterns(t *testing.T) {
	cfg := &config.Config{
		DefaultBrowser: "chrome",
		Rules: []config.Rule{
			{
				Match:   []string{"*.google.com", "*.youtube.com"},
				Browser: "firefox",
			},
		},
	}

	router := NewRouter(cfg)

	tests := []struct {
		name        string
		url         string
		wantBrowser string
		wantMatched bool
	}{
		{
			name:        "matches first pattern",
			url:         "https://mail.google.com",
			wantBrowser: "firefox",
			wantMatched: true,
		},
		{
			name:        "matches second pattern",
			url:         "https://www.youtube.com",
			wantBrowser: "firefox",
			wantMatched: true,
		},
		{
			name:        "no match",
			url:         "https://example.com",
			wantBrowser: "chrome",
			wantMatched: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			browser, _, matched := router.FindMatch(tt.url)
			if browser != tt.wantBrowser {
				t.Errorf("FindMatch() browser = %q, want %q", browser, tt.wantBrowser)
			}
			if matched != tt.wantMatched {
				t.Errorf("FindMatch() matched = %v, want %v", matched, tt.wantMatched)
			}
		})
	}
}
