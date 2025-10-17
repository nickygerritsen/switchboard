package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nickygerritsen/switchboard/internal/browser"
	"github.com/nickygerritsen/switchboard/internal/config"
	"github.com/spf13/cobra"
)

func TestRunListBrowsers(t *testing.T) {
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
rules: []
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Save the original factories and cfgFile
	originalDetectorFactory := detectorFactory
	originalCfgFile := cfgFile
	defer func() {
		detectorFactory = originalDetectorFactory
		cfgFile = originalCfgFile
	}()

	cfgFile = configPath

	tests := []struct {
		name        string
		setupFake   func() *fakeDetector
		wantErr     bool
		wantOutput  []string
		wantContain string
	}{
		{
			name: "multiple browsers detected",
			setupFake: func() *fakeDetector {
				return &fakeDetector{
					browsers: map[string]*browser.Browser{
						"firefox": {Name: "firefox", Path: "/usr/bin/firefox"},
						"chrome":  {Name: "chrome", Path: "/usr/bin/google-chrome"},
						"safari":  {Name: "safari", Path: "/Applications/Safari.app/Contents/MacOS/Safari"},
					},
				}
			},
			wantErr: false,
			wantOutput: []string{
				"Detected 3 browser(s):",
				"firefox",
				"/usr/bin/firefox",
				"chrome",
				"/usr/bin/google-chrome",
				"safari",
				"/Applications/Safari.app/Contents/MacOS/Safari",
			},
		},
		{
			name: "single browser detected",
			setupFake: func() *fakeDetector {
				return &fakeDetector{
					browsers: map[string]*browser.Browser{
						"firefox": {Name: "firefox", Path: "/usr/bin/firefox"},
					},
				}
			},
			wantErr: false,
			wantOutput: []string{
				"Detected 1 browser(s):",
				"firefox",
				"/usr/bin/firefox",
			},
		},
		{
			name: "no browsers detected",
			setupFake: func() *fakeDetector {
				return &fakeDetector{
					browsers: map[string]*browser.Browser{},
				}
			},
			wantErr:     false,
			wantContain: "No browsers detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := tt.setupFake()

			// Override detector factory
			detectorFactory = func(cfg *config.Config) browserDetector {
				return detector
			}

			// Capture output
			var buf bytes.Buffer
			listCmd := &cobra.Command{
				Use:  listBrowsersCmd.Use,
				Args: listBrowsersCmd.Args,
				RunE: listBrowsersCmd.RunE,
			}
			listCmd.SetOut(&buf)
			listCmd.SetErr(&buf)

			err := runListBrowsers(listCmd, []string{})

			if (err != nil) != tt.wantErr {
				t.Errorf("runListBrowsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			output := buf.String()

			// Check expected output
			for _, want := range tt.wantOutput {
				if !strings.Contains(output, want) {
					t.Errorf("runListBrowsers() output missing %q, got:\n%s", want, output)
				}
			}

			// Check for specific string
			if tt.wantContain != "" {
				if !strings.Contains(output, tt.wantContain) {
					t.Errorf("runListBrowsers() output should contain %q, got:\n%s", tt.wantContain, output)
				}
			}
		})
	}
}

func TestListBrowsersCmdValidation(t *testing.T) {
	// list-browsers command takes no arguments
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no arguments is valid",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "with arguments is valid (cobra doesn't restrict by default)",
			args:    []string{"extra"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// list-browsers doesn't have custom Args validation, so we just test it doesn't panic
			if listBrowsersCmd.Args != nil {
				err := listBrowsersCmd.Args(listBrowsersCmd, tt.args)
				if (err != nil) != tt.wantErr {
					t.Errorf("listBrowsersCmd.Args() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestListBrowsersCmdMetadata(t *testing.T) {
	if listBrowsersCmd.Use != "list-browsers" {
		t.Errorf("listBrowsersCmd.Use = %q, want %q", listBrowsersCmd.Use, "list-browsers")
	}

	if listBrowsersCmd.Short == "" {
		t.Error("listBrowsersCmd.Short should not be empty")
	}

	if listBrowsersCmd.Long == "" {
		t.Error("listBrowsersCmd.Long should not be empty")
	}

	if listBrowsersCmd.RunE == nil {
		t.Error("listBrowsersCmd.RunE should not be nil")
	}
}
