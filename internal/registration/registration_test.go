package registration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetBinaryPath(t *testing.T) {
	path, err := getBinaryPath()
	if err != nil {
		t.Fatalf("getBinaryPath() returned error: %v", err)
	}

	if path == "" {
		t.Error("getBinaryPath() returned empty path")
	}

	// Should be an absolute path
	if !filepath.IsAbs(path) {
		t.Errorf("getBinaryPath() returned relative path: %s", path)
	}

	// Should exist
	if _, err := os.Stat(path); err != nil {
		t.Errorf("getBinaryPath() returned non-existent path: %s, error: %v", path, err)
	}
}

func TestNew(t *testing.T) {
	registrar, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if registrar == nil {
		t.Fatal("New() returned nil registrar")
	}
}

func TestRegistrarGetBinaryPath(t *testing.T) {
	registrar, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	binPath, err := registrar.GetBinaryPath()
	if err != nil {
		t.Errorf("GetBinaryPath() returned error: %v", err)
	}
	if binPath == "" {
		t.Error("GetBinaryPath() returned empty path")
	}
	if !filepath.IsAbs(binPath) {
		t.Errorf("GetBinaryPath() returned relative path: %s", binPath)
	}
}

func TestRegistrarIsRegistered(t *testing.T) {
	registrar, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	// IsRegistered should not error (even if not registered)
	_, err = registrar.IsRegistered()
	if err != nil {
		t.Errorf("IsRegistered() returned error: %v", err)
	}
}
