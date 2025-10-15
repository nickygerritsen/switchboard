package browser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nickygerritsen/switchboard/internal/config"
)

func TestGetBrowserDef(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantNil  bool
	}{
		{
			name:     "chrome by name",
			input:    "chrome",
			wantName: "chrome",
			wantNil:  false,
		},
		{
			name:     "firefox by name",
			input:    "firefox",
			wantName: "firefox",
			wantNil:  false,
		},
		{
			name:     "brave by name",
			input:    "brave",
			wantName: "brave",
			wantNil:  false,
		},
		{
			name:     "safari by name",
			input:    "safari",
			wantName: "safari",
			wantNil:  false,
		},
		{
			name:    "unknown browser",
			input:   "unknown-browser",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBrowserDef(tt.input)
			if tt.wantNil {
				if result != nil {
					t.Errorf("GetBrowserDef() = %v, want nil", result)
				}
			} else {
				if result == nil {
					t.Fatalf("GetBrowserDef() = nil, want browser def")
				}
				if result.Name != tt.wantName {
					t.Errorf("GetBrowserDef().Name = %v, want %v", result.Name, tt.wantName)
				}
			}
		})
	}
}

func TestGetAllBrowserDefs(t *testing.T) {
	defs := GetAllBrowserDefs()
	if len(defs) == 0 {
		t.Error("GetAllBrowserDefs() returned empty list")
	}

	// Check that essential browsers are included
	names := make(map[string]bool)
	for _, def := range defs {
		names[def.Name] = true
	}

	essential := []string{"chrome", "firefox", "brave", "safari"}
	for _, name := range essential {
		if !names[name] {
			t.Errorf("GetAllBrowserDefs() missing essential browser: %s", name)
		}
	}
}

func TestNewDetector(t *testing.T) {
	cfg := &config.Config{
		DefaultBrowser: "chrome",
	}

	detector := NewDetector(cfg)
	if detector == nil {
		t.Fatal("NewDetector() returned nil")
	}
	if detector.config != cfg {
		t.Error("NewDetector() did not store config")
	}
	if detector.cache == nil {
		t.Error("NewDetector() did not initialize cache")
	}
}

func TestDetector_DetectWithCustomPath(t *testing.T) {
	// Create a temporary file to act as a browser
	tmpDir := t.TempDir()
	browserPath := filepath.Join(tmpDir, "fake-browser")
	f, err := os.Create(browserPath)
	if err != nil {
		t.Fatalf("Failed to create fake browser: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Failed to close fake browser: %v", err)
	}

	cfg := &config.Config{
		DefaultBrowser: "chrome",
		Browsers: map[string]config.Browser{
			"chrome": {
				Path: browserPath,
			},
		},
	}

	detector := NewDetector(cfg)
	browser, err := detector.Detect("chrome")
	if err != nil {
		t.Fatalf("Detect() failed: %v", err)
	}
	if browser.Path != browserPath {
		t.Errorf("Detect() path = %v, want %v", browser.Path, browserPath)
	}
	if browser.Name != "chrome" {
		t.Errorf("Detect() name = %v, want chrome", browser.Name)
	}
}

func TestDetector_DetectUnknownBrowser(t *testing.T) {
	cfg := &config.Config{
		DefaultBrowser: "chrome",
	}

	detector := NewDetector(cfg)
	_, err := detector.Detect("unknown-browser")
	if err == nil {
		t.Error("Detect() should fail for unknown browser")
	}
}

func TestDetector_Cache(t *testing.T) {
	// Create a temporary file to act as a browser
	tmpDir := t.TempDir()
	browserPath := filepath.Join(tmpDir, "fake-browser")
	f, err := os.Create(browserPath)
	if err != nil {
		t.Fatalf("Failed to create fake browser: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Failed to close fake browser: %v", err)
	}

	cfg := &config.Config{
		DefaultBrowser: "chrome",
		Browsers: map[string]config.Browser{
			"chrome": {
				Path: browserPath,
			},
		},
	}

	detector := NewDetector(cfg)

	// First detection should populate cache
	browser1, err := detector.Detect("chrome")
	if err != nil {
		t.Fatalf("First Detect() failed: %v", err)
	}

	// Check cache
	cached, ok := detector.GetCached("chrome")
	if !ok {
		t.Error("Browser should be in cache after first detection")
	}
	if cached != browser1 {
		t.Error("Cached browser should be the same instance")
	}

	// Second detection should use cache
	browser2, err := detector.Detect("chrome")
	if err != nil {
		t.Fatalf("Second Detect() failed: %v", err)
	}
	if browser2 != browser1 {
		t.Error("Second detection should return cached browser")
	}
}

func TestDetector_ClearCache(t *testing.T) {
	tmpDir := t.TempDir()
	browserPath := filepath.Join(tmpDir, "fake-browser")
	f, err := os.Create(browserPath)
	if err != nil {
		t.Fatalf("Failed to create fake browser: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Failed to close fake browser: %v", err)
	}

	cfg := &config.Config{
		DefaultBrowser: "chrome",
		Browsers: map[string]config.Browser{
			"chrome": {
				Path: browserPath,
			},
		},
	}

	detector := NewDetector(cfg)

	// Populate cache
	_, err = detector.Detect("chrome")
	if err != nil {
		t.Fatalf("Detect() failed: %v", err)
	}

	// Verify cache has entry
	_, ok := detector.GetCached("chrome")
	if !ok {
		t.Error("Cache should have entry before clear")
	}

	// Clear cache
	detector.ClearCache()

	// Verify cache is empty
	_, ok = detector.GetCached("chrome")
	if ok {
		t.Error("Cache should be empty after clear")
	}
}

func TestDetector_DetectAll(t *testing.T) {
	cfg := &config.Config{
		DefaultBrowser: "chrome",
	}

	detector := NewDetector(cfg)
	browsers := detector.DetectAll()

	// The result should be a map (may be empty if no browsers are installed)
	if browsers == nil {
		t.Error("DetectAll() should return non-nil map")
	}

	// On most development machines, at least one browser should be found
	// But we can't guarantee this in CI, so we just check the map is valid
	for name, browser := range browsers {
		if browser == nil {
			t.Errorf("DetectAll() returned nil browser for %s", name)
			continue
		}
		if browser.Name != name {
			t.Errorf("DetectAll() browser name mismatch: got %s, want %s", browser.Name, name)
		}
		if browser.Path == "" {
			t.Errorf("DetectAll() browser %s has empty path", name)
		}
	}
}

func TestBrowserDef_GetPaths(t *testing.T) {
	chrome := GetBrowserDef("chrome")
	if chrome == nil {
		t.Fatal("GetBrowserDef('chrome') returned nil")
	}

	paths := chrome.GetPaths()
	if len(paths) == 0 {
		t.Error("GetPaths() returned empty list for chrome")
	}

	// Paths should be OS-specific
	// We can't test exact paths as they vary by OS, but we can check they're not empty
	for _, path := range paths {
		if path == "" {
			t.Error("GetPaths() returned empty path")
		}
	}
}
