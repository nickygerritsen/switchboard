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

func TestRunInit(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(tempDir string) string // Returns config path to use
		wantErr        bool
		wantErrContain string
		wantOutput     []string
		checkFile      func(t *testing.T, configPath string)
	}{
		{
			name: "successful init - new config",
			setupFunc: func(tempDir string) string {
				return filepath.Join(tempDir, "config.yaml")
			},
			wantErr: false,
			wantOutput: []string{
				"Created configuration file at:",
				"Edit this file to customize your browser routing rules",
			},
			checkFile: func(t *testing.T, configPath string) {
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					t.Errorf("runInit() did not create config file at %s", configPath)
				} else {
					// Verify the file has valid content
					data, err := os.ReadFile(configPath)
					if err != nil {
						t.Errorf("Failed to read created config: %v", err)
					}
					content := string(data)
					// Check for some expected content
					if !strings.Contains(content, "defaultBrowser") {
						t.Error("Created config should contain defaultBrowser field")
					}
				}
			},
		},
		{
			name: "init fails - config already exists",
			setupFunc: func(tempDir string) string {
				configPath := filepath.Join(tempDir, "config.yaml")
				_ = os.WriteFile(configPath, []byte("defaultBrowser: firefox\n"), 0644)
				return configPath
			},
			wantErr:        true,
			wantErrContain: "configuration file already exists",
		},
		{
			name: "init fails - invalid directory",
			setupFunc: func(tempDir string) string {
				// Return a path in a non-existent directory
				return "/nonexistent/directory/that/does/not/exist/config.yaml"
			},
			wantErr:        true,
			wantErrContain: "failed to save config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configPath := tt.setupFunc(tempDir)

			// Use dependency injection to override config path
			oldProvider := config.SetConfigPathProvider(func() (string, error) {
				return configPath, nil
			})
			defer config.SetConfigPathProvider(oldProvider)

			// Capture output
			var buf bytes.Buffer
			iniCmd := &cobra.Command{
				Use:  initCmd.Use,
				Args: initCmd.Args,
				RunE: initCmd.RunE,
			}
			iniCmd.SetOut(&buf)
			iniCmd.SetErr(&buf)

			err := runInit(iniCmd, []string{})

			if (err != nil) != tt.wantErr {
				t.Errorf("runInit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.wantErrContain != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContain) {
					t.Errorf("runInit() error should contain %q, got %v", tt.wantErrContain, err)
				}
			}

			if !tt.wantErr {
				output := buf.String()
				for _, want := range tt.wantOutput {
					if !strings.Contains(output, want) {
						t.Errorf("runInit() output missing %q, got:\n%s", want, output)
					}
				}

				// Check file if needed
				if tt.checkFile != nil {
					tt.checkFile(t, configPath)
				}
			}
		})
	}
}

func TestRunInitGetConfigPathError(t *testing.T) {
	// Use dependency injection to return an error
	oldProvider := config.SetConfigPathProvider(func() (string, error) {
		return "", os.ErrPermission
	})
	defer config.SetConfigPathProvider(oldProvider)

	err := runInit(initCmd, []string{})
	if err == nil {
		t.Error("runInit() should fail when GetConfigPath returns an error")
	}
	if !strings.Contains(err.Error(), "failed to get config path") {
		t.Errorf("runInit() error should mention config path, got: %v", err)
	}
}

func TestInitCmdMetadata(t *testing.T) {
	if initCmd.Use != "init" {
		t.Errorf("initCmd.Use = %q, want %q", initCmd.Use, "init")
	}

	if initCmd.Short == "" {
		t.Error("initCmd.Short should not be empty")
	}

	if initCmd.Long == "" {
		t.Error("initCmd.Long should not be empty")
	}

	if initCmd.RunE == nil {
		t.Error("initCmd.RunE should not be nil")
	}
}
