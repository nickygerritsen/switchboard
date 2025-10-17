package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunUnregister(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func() browserRegistrar
		wantErr        bool
		wantErrContain string
		wantOutput     []string
	}{
		{
			name: "successful unregistration - currently registered",
			setupMock: func() browserRegistrar {
				return &mockRegistrar{
					isRegisteredFunc: func() (bool, error) {
						return true, nil
					},
					unregisterFunc: func() error {
						return nil
					},
				}
			},
			wantErr: false,
			wantOutput: []string{
				"Successfully unregistered Switchboard as a browser",
			},
		},
		{
			name: "not registered",
			setupMock: func() browserRegistrar {
				return &mockRegistrar{
					isRegisteredFunc: func() (bool, error) {
						return false, nil
					},
				}
			},
			wantErr: false,
			wantOutput: []string{
				"Switchboard is not currently registered as a browser",
			},
		},
		{
			name: "IsRegistered check fails",
			setupMock: func() browserRegistrar {
				return &mockRegistrar{
					isRegisteredFunc: func() (bool, error) {
						return false, errors.New("permission denied")
					},
				}
			},
			wantErr:        true,
			wantErrContain: "failed to check registration status",
		},
		{
			name: "Unregister operation fails",
			setupMock: func() browserRegistrar {
				return &mockRegistrar{
					isRegisteredFunc: func() (bool, error) {
						return true, nil
					},
					unregisterFunc: func() error {
						return errors.New("registry access denied")
					},
				}
			},
			wantErr:        true,
			wantErrContain: "failed to unregister",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override factory
			oldFactory := registrarFactory
			registrarFactory = func() (browserRegistrar, error) {
				return tt.setupMock(), nil
			}
			defer func() { registrarFactory = oldFactory }()

			// Capture output
			var buf bytes.Buffer
			cmd := &cobra.Command{
				Use:  unregisterCmd.Use,
				RunE: unregisterCmd.RunE,
			}
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			err := runUnregister(cmd, []string{})

			if (err != nil) != tt.wantErr {
				t.Errorf("runUnregister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.wantErrContain != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContain) {
					t.Errorf("runUnregister() error should contain %q, got %v", tt.wantErrContain, err)
				}
			}

			if !tt.wantErr {
				output := buf.String()
				for _, want := range tt.wantOutput {
					if !strings.Contains(output, want) {
						t.Errorf("runUnregister() output missing %q, got:\n%s", want, output)
					}
				}
			}
		})
	}
}

func TestRunUnregisterFactoryError(t *testing.T) {
	// Override factory to return error
	oldFactory := registrarFactory
	registrarFactory = func() (browserRegistrar, error) {
		return nil, errors.New("unsupported platform")
	}
	defer func() { registrarFactory = oldFactory }()

	err := runUnregister(unregisterCmd, []string{})
	if err == nil {
		t.Error("runUnregister() should fail when registrarFactory returns an error")
	}
	if !strings.Contains(err.Error(), "failed to create registrar") {
		t.Errorf("runUnregister() error should mention registrar creation, got: %v", err)
	}
}

func TestUnregisterCmdMetadata(t *testing.T) {
	if unregisterCmd.Use != "unregister" {
		t.Errorf("unregisterCmd.Use = %q, want %q", unregisterCmd.Use, "unregister")
	}

	if unregisterCmd.Short == "" {
		t.Error("unregisterCmd.Short should not be empty")
	}

	if unregisterCmd.Long == "" {
		t.Error("unregisterCmd.Long should not be empty")
	}

	if unregisterCmd.RunE == nil {
		t.Error("unregisterCmd.RunE should not be nil")
	}
}
