//go:build windows

package registration

import (
	"path/filepath"
	"testing"
)

func TestNewWindowsRegistrar(t *testing.T) {
	reg, err := newWindowsRegistrar()
	if err != nil {
		t.Fatalf("newWindowsRegistrar() returned error: %v", err)
	}

	if reg == nil {
		t.Fatal("newWindowsRegistrar() returned nil")
	}

	if reg.binaryPath == "" {
		t.Error("newWindowsRegistrar() created registrar with empty binaryPath")
	}
}

func TestWindowsRegistrarGetBinaryPath(t *testing.T) {
	reg, err := newWindowsRegistrar()
	if err != nil {
		t.Fatalf("newWindowsRegistrar() returned error: %v", err)
	}

	path, err := reg.GetBinaryPath()
	if err != nil {
		t.Fatalf("GetBinaryPath() returned error: %v", err)
	}

	if path == "" {
		t.Error("GetBinaryPath() returned empty path")
	}

	if !filepath.IsAbs(path) {
		t.Errorf("GetBinaryPath() returned relative path: %s", path)
	}
}

func TestWindowsRegistrarIsRegistered(t *testing.T) {
	reg, err := newWindowsRegistrar()
	if err != nil {
		t.Fatalf("newWindowsRegistrar() returned error: %v", err)
	}

	// IsRegistered should not error (even if not registered)
	_, err = reg.IsRegistered()
	if err != nil {
		t.Errorf("IsRegistered() returned error: %v", err)
	}
}
