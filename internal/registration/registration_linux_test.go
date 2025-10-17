//go:build linux

package registration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewLinuxRegistrar(t *testing.T) {
	reg, err := newLinuxRegistrar()
	if err != nil {
		t.Fatalf("newLinuxRegistrar() returned error: %v", err)
	}

	if reg == nil {
		t.Fatal("newLinuxRegistrar() returned nil")
	}

	if reg.binaryPath == "" {
		t.Error("newLinuxRegistrar() created registrar with empty binaryPath")
	}
}

func TestLinuxRegistrarGetBinaryPath(t *testing.T) {
	reg, err := newLinuxRegistrar()
	if err != nil {
		t.Fatalf("newLinuxRegistrar() returned error: %v", err)
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

func TestLinuxRegistrarGetDesktopFilePath(t *testing.T) {
	reg := &linuxRegistrar{binaryPath: "/usr/local/bin/switchboard"}

	desktopPath := reg.getDesktopFilePath()

	if desktopPath == "" {
		t.Error("getDesktopFilePath() returned empty path")
	}

	if !strings.HasSuffix(desktopPath, "switchboard.desktop") {
		t.Errorf("getDesktopFilePath() should end with switchboard.desktop, got: %s", desktopPath)
	}

	if !strings.Contains(desktopPath, ".local/share/applications") {
		t.Errorf("getDesktopFilePath() should contain .local/share/applications, got: %s", desktopPath)
	}

	home, _ := os.UserHomeDir()
	expectedPath := filepath.Join(home, ".local", "share", "applications", "switchboard.desktop")
	if desktopPath != expectedPath {
		t.Errorf("getDesktopFilePath() = %s, want %s", desktopPath, expectedPath)
	}
}

func TestLinuxRegistrarGenerateDesktopFile(t *testing.T) {
	reg := &linuxRegistrar{binaryPath: "/usr/local/bin/switchboard"}

	desktop := reg.generateDesktopFile()

	if desktop == "" {
		t.Error("generateDesktopFile() returned empty string")
	}

	// Check for required entries
	requiredEntries := []string{
		"[Desktop Entry]",
		"Version=1.0",
		"Type=Application",
		"Name=Switchboard",
		"Exec=/usr/local/bin/switchboard open %u",
		"Terminal=false",
		"Categories=Network;WebBrowser;",
		"MimeType=x-scheme-handler/http;x-scheme-handler/https;text/html;",
	}

	for _, entry := range requiredEntries {
		if !strings.Contains(desktop, entry) {
			t.Errorf("generateDesktopFile() missing required entry: %s", entry)
		}
	}

	// Check it starts with [Desktop Entry]
	if !strings.HasPrefix(desktop, "[Desktop Entry]") {
		t.Error("generateDesktopFile() should start with [Desktop Entry]")
	}
}
