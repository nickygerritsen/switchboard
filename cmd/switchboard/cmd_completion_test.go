package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunCompletion(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantErr      bool
		wantContains string
	}{
		{
			name:         "bash completion",
			args:         []string{"bash"},
			wantErr:      false,
			wantContains: "# bash completion for switchboard",
		},
		{
			name:         "zsh completion",
			args:         []string{"zsh"},
			wantErr:      false,
			wantContains: "#compdef switchboard",
		},
		{
			name:         "fish completion",
			args:         []string{"fish"},
			wantErr:      false,
			wantContains: "# fish completion for switchboard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outBuf bytes.Buffer
			cmd := &cobra.Command{
				Use:  completionCmd.Use,
				RunE: completionCmd.RunE,
			}
			cmd.SetOut(&outBuf)
			cmd.SetErr(&outBuf)

			err := runCompletion(cmd, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			output := outBuf.String()
			if !strings.Contains(output, tt.wantContains) {
				t.Errorf("Expected output to contain %q, but it didn't.\nOutput first 100 chars: %s", tt.wantContains, truncate(output, 100))
			}
		})
	}
}

func TestCompletionCmdMetadata(t *testing.T) {
	if completionCmd.Use != "completion [bash|zsh|fish]" {
		t.Errorf("completionCmd.Use = %q, want %q", completionCmd.Use, "completion [bash|zsh|fish]")
	}

	if completionCmd.Short == "" {
		t.Error("completionCmd.Short should not be empty")
	}

	if completionCmd.Long == "" {
		t.Error("completionCmd.Long should not be empty")
	}

	if completionCmd.RunE == nil {
		t.Error("completionCmd.RunE should not be nil")
	}

	// Check valid args - should contain bash, zsh, and fish
	expectedValidArgs := map[string]bool{"bash": true, "zsh": true, "fish": true}
	if len(completionCmd.ValidArgs) != len(expectedValidArgs) {
		t.Errorf("completionCmd.ValidArgs length = %d, want %d", len(completionCmd.ValidArgs), len(expectedValidArgs))
	}
	for _, arg := range completionCmd.ValidArgs {
		if !expectedValidArgs[arg] {
			t.Errorf("Unexpected valid arg: %q", arg)
		}
	}
	for arg := range expectedValidArgs {
		found := false
		for _, validArg := range completionCmd.ValidArgs {
			if validArg == arg {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected valid arg %q not found in ValidArgs", arg)
		}
	}

	// Check that help text contains key information
	if !strings.Contains(completionCmd.Long, "bash") {
		t.Error("Long help should mention bash")
	}
	if !strings.Contains(completionCmd.Long, "zsh") {
		t.Error("Long help should mention zsh")
	}
	if !strings.Contains(completionCmd.Long, "fish") {
		t.Error("Long help should mention fish")
	}
}

// truncate returns the first n characters of s
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
