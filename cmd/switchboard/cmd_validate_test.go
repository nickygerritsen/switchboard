package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunValidate(t *testing.T) {
	// Save original cfgFile and restore after test
	originalCfgFile := cfgFile
	defer func() { cfgFile = originalCfgFile }()

	tests := []struct {
		name           string
		configContent  string
		wantErr        bool
		wantErrContain string
		wantOutput     []string
	}{
		{
			name: "valid configuration",
			configContent: `defaultBrowser: firefox
debug: true
logFile: /tmp/test.log
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
`,
			wantErr: false,
			wantOutput: []string{
				"Configuration is valid",
				"Default browser: firefox",
				"Rules: 2",
				"Debug: true",
				"Log file: /tmp/test.log",
			},
		},
		{
			name: "valid configuration without debug",
			configContent: `defaultBrowser: chrome
debug: false
logFile: ""
browsers:
  chrome:
    path: /usr/bin/chrome
rules: []
`,
			wantErr: false,
			wantOutput: []string{
				"Configuration is valid",
				"Default browser: chrome",
				"Rules: 0",
				"Debug: false",
				"Log file:",
				"(default)",
			},
		},
		{
			name: "invalid configuration - missing defaultBrowser",
			configContent: `debug: false
logFile: ""
browsers:
  chrome:
    path: /usr/bin/chrome
rules: []
`,
			wantErr:        true,
			wantErrContain: "configuration validation failed",
		},
		{
			name: "invalid configuration - invalid YAML",
			configContent: `defaultBrowser: firefox
debug: [this is not valid
browsers:
`,
			wantErr:        true,
			wantErrContain: "configuration validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "config.yaml")

			if err := os.WriteFile(configPath, []byte(tt.configContent), 0644); err != nil {
				t.Fatalf("Failed to create test config: %v", err)
			}

			cfgFile = configPath

			// Capture output
			var buf bytes.Buffer
			valCmd := &cobra.Command{
				Use:  validateCmd.Use,
				Args: validateCmd.Args,
				RunE: validateCmd.RunE,
			}
			valCmd.SetOut(&buf)
			valCmd.SetErr(&buf)

			err := runValidate(valCmd, []string{})

			if (err != nil) != tt.wantErr {
				t.Errorf("runValidate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.wantErrContain != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContain) {
					t.Errorf("runValidate() error should contain %q, got %v", tt.wantErrContain, err)
				}
			}

			if !tt.wantErr {
				output := buf.String()
				for _, want := range tt.wantOutput {
					if !strings.Contains(output, want) {
						t.Errorf("runValidate() output missing %q, got:\n%s", want, output)
					}
				}
			}
		})
	}
}

func TestRunValidateNonexistentFile(t *testing.T) {
	// Save original cfgFile and restore after test
	originalCfgFile := cfgFile
	defer func() { cfgFile = originalCfgFile }()

	cfgFile = "/nonexistent/path/to/config.yaml"

	err := runValidate(validateCmd, []string{})
	if err == nil {
		t.Error("runValidate() should fail with nonexistent config file")
	}
}

func TestValidateCmdMetadata(t *testing.T) {
	if validateCmd.Use != "validate" {
		t.Errorf("validateCmd.Use = %q, want %q", validateCmd.Use, "validate")
	}

	if validateCmd.Short == "" {
		t.Error("validateCmd.Short should not be empty")
	}

	if validateCmd.Long == "" {
		t.Error("validateCmd.Long should not be empty")
	}

	if validateCmd.RunE == nil {
		t.Error("validateCmd.RunE should not be nil")
	}
}
