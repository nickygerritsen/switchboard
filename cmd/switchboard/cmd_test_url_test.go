package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nickygerritsen/switchboard/internal/config"
	"github.com/spf13/cobra"
)

func TestRunTest(t *testing.T) {
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
    profile: Work
  - match:
      - "*.google.com"
    browser: chrome
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Save the original factories and cfgFile
	originalRouterFactory := routerFactory
	originalCfgFile := cfgFile
	defer func() {
		routerFactory = originalRouterFactory
		cfgFile = originalCfgFile
	}()

	cfgFile = configPath

	tests := []struct {
		name           string
		args           []string
		setupRouter    func() *fakeRouter
		wantErr        bool
		wantOutput     []string
		wantNotContain []string
	}{
		{
			name: "matched URL with profile",
			args: []string{"https://github.com/test"},
			setupRouter: func() *fakeRouter {
				return &fakeRouter{
					matches: map[string]routeResult{
						"https://github.com/test": {browser: "firefox", profile: "Work", matched: true},
					},
				}
			},
			wantErr: false,
			wantOutput: []string{
				"URL: https://github.com/test",
				"Browser: firefox",
				"Profile: Work",
				"Matched: yes",
			},
		},
		{
			name: "matched URL without profile",
			args: []string{"https://google.com/search"},
			setupRouter: func() *fakeRouter {
				return &fakeRouter{
					matches: map[string]routeResult{
						"https://google.com/search": {browser: "chrome", profile: "", matched: true},
					},
				}
			},
			wantErr: false,
			wantOutput: []string{
				"URL: https://google.com/search",
				"Browser: chrome",
				"Matched: yes",
			},
			wantNotContain: []string{
				"Profile:",
			},
		},
		{
			name: "unmatched URL uses default",
			args: []string{"https://example.com"},
			setupRouter: func() *fakeRouter {
				return &fakeRouter{
					matches: map[string]routeResult{
						"https://example.com": {browser: "firefox", profile: "", matched: false},
					},
				}
			},
			wantErr: false,
			wantOutput: []string{
				"URL: https://example.com",
				"Browser: firefox (default)",
				"Matched: no",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := tt.setupRouter()

			// Override router factory
			routerFactory = func(cfg *config.Config) urlRouter {
				return router
			}

			// Capture output
			var buf bytes.Buffer
			testCmd := &cobra.Command{
				Use:  testCmd.Use,
				Args: testCmd.Args,
				RunE: testCmd.RunE,
			}
			testCmd.SetOut(&buf)
			testCmd.SetErr(&buf)

			err := runTest(testCmd, tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("runTest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			output := buf.String()

			// Check expected output
			for _, want := range tt.wantOutput {
				if !strings.Contains(output, want) {
					t.Errorf("runTest() output missing %q, got:\n%s", want, output)
				}
			}

			// Check for strings that should not be present
			for _, notWant := range tt.wantNotContain {
				if strings.Contains(output, notWant) {
					t.Errorf("runTest() output should not contain %q, got:\n%s", notWant, output)
				}
			}
		})
	}
}

func TestTestCmdValidation(t *testing.T) {
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
			err := testCmd.Args(testCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("testCmd.Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTestCmdMetadata(t *testing.T) {
	if testCmd.Use != "test <url>" {
		t.Errorf("testCmd.Use = %q, want %q", testCmd.Use, "test <url>")
	}

	if testCmd.Short == "" {
		t.Error("testCmd.Short should not be empty")
	}

	if testCmd.Long == "" {
		t.Error("testCmd.Long should not be empty")
	}

	if testCmd.RunE == nil {
		t.Error("testCmd.RunE should not be nil")
	}
}
