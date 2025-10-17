package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// mockRegistrar is a mock implementation of browserRegistrar for testing
type mockRegistrar struct {
	registerFunc    func() error
	unregisterFunc  func() error
	isRegisteredFunc func() (bool, error)
	getBinaryPathFunc func() (string, error)
}

func (m *mockRegistrar) Register() error {
	if m.registerFunc != nil {
		return m.registerFunc()
	}
	return nil
}

func (m *mockRegistrar) Unregister() error {
	if m.unregisterFunc != nil {
		return m.unregisterFunc()
	}
	return nil
}

func (m *mockRegistrar) IsRegistered() (bool, error) {
	if m.isRegisteredFunc != nil {
		return m.isRegisteredFunc()
	}
	return false, nil
}

func (m *mockRegistrar) GetBinaryPath() (string, error) {
	if m.getBinaryPathFunc != nil {
		return m.getBinaryPathFunc()
	}
	return "/usr/local/bin/switchboard", nil
}

func TestRunRegister(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func() browserRegistrar
		wantErr        bool
		wantErrContain string
		wantOutput     []string
	}{
		{
			name: "successful registration - not yet registered",
			setupMock: func() browserRegistrar {
				return &mockRegistrar{
					isRegisteredFunc: func() (bool, error) {
						return false, nil
					},
					registerFunc: func() error {
						return nil
					},
				}
			},
			wantErr: false,
			wantOutput: []string{
				"Successfully registered Switchboard as a browser!",
				"To set Switchboard as your default browser",
			},
		},
		{
			name: "already registered",
			setupMock: func() browserRegistrar {
				return &mockRegistrar{
					isRegisteredFunc: func() (bool, error) {
						return true, nil
					},
				}
			},
			wantErr: false,
			wantOutput: []string{
				"Switchboard is already registered as a browser",
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
			name: "Register operation fails",
			setupMock: func() browserRegistrar {
				return &mockRegistrar{
					isRegisteredFunc: func() (bool, error) {
						return false, nil
					},
					registerFunc: func() error {
						return errors.New("registry access denied")
					},
				}
			},
			wantErr:        true,
			wantErrContain: "failed to register",
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
				Use:  registerCmd.Use,
				RunE: registerCmd.RunE,
			}
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			err := runRegister(cmd, []string{})

			if (err != nil) != tt.wantErr {
				t.Errorf("runRegister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.wantErrContain != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContain) {
					t.Errorf("runRegister() error should contain %q, got %v", tt.wantErrContain, err)
				}
			}

			if !tt.wantErr {
				output := buf.String()
				for _, want := range tt.wantOutput {
					if !strings.Contains(output, want) {
						t.Errorf("runRegister() output missing %q, got:\n%s", want, output)
					}
				}
			}
		})
	}
}

func TestRunRegisterFactoryError(t *testing.T) {
	// Override factory to return error
	oldFactory := registrarFactory
	registrarFactory = func() (browserRegistrar, error) {
		return nil, errors.New("unsupported platform")
	}
	defer func() { registrarFactory = oldFactory }()

	err := runRegister(registerCmd, []string{})
	if err == nil {
		t.Error("runRegister() should fail when registrarFactory returns an error")
	}
	if !strings.Contains(err.Error(), "failed to create registrar") {
		t.Errorf("runRegister() error should mention registrar creation, got: %v", err)
	}
}

func TestRegisterCmdMetadata(t *testing.T) {
	if registerCmd.Use != "register" {
		t.Errorf("registerCmd.Use = %q, want %q", registerCmd.Use, "register")
	}

	if registerCmd.Short == "" {
		t.Error("registerCmd.Short should not be empty")
	}

	if registerCmd.Long == "" {
		t.Error("registerCmd.Long should not be empty")
	}

	if registerCmd.RunE == nil {
		t.Error("registerCmd.RunE should not be nil")
	}
}
