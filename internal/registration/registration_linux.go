//go:build linux

package registration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// New creates a new Registrar for the current platform
func New() (Registrar, error) {
	return newLinuxRegistrar()
}

type linuxRegistrar struct {
	binaryPath string
}

func newLinuxRegistrar() (*linuxRegistrar, error) {
	binPath, err := getBinaryPath()
	if err != nil {
		return nil, err
	}

	return &linuxRegistrar{
		binaryPath: binPath,
	}, nil
}

func (r *linuxRegistrar) Register() error {
	// Create .desktop file
	if err := r.createDesktopFile(); err != nil {
		return fmt.Errorf("failed to create desktop file: %w", err)
	}

	// Update desktop database
	if err := r.updateDesktopDatabase(); err != nil {
		return fmt.Errorf("failed to update desktop database: %w", err)
	}

	// Try to set as default browser (requires xdg-settings)
	_ = r.setAsDefaultBrowser() // Don't fail if this doesn't work

	return nil
}

func (r *linuxRegistrar) Unregister() error {
	desktopPath := r.getDesktopFilePath()

	// Remove desktop file
	if err := os.Remove(desktopPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove desktop file: %w", err)
	}

	// Update desktop database
	return r.updateDesktopDatabase()
}

func (r *linuxRegistrar) IsRegistered() (bool, error) {
	desktopPath := r.getDesktopFilePath()

	// Check if desktop file exists
	_, err := os.Stat(desktopPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *linuxRegistrar) GetBinaryPath() (string, error) {
	return r.binaryPath, nil
}

// getDesktopFilePath returns the path to the .desktop file
func (r *linuxRegistrar) getDesktopFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "applications", "switchboard.desktop")
}

// createDesktopFile creates the .desktop file for Switchboard
func (r *linuxRegistrar) createDesktopFile() error {
	desktopPath := r.getDesktopFilePath()

	// Ensure directory exists
	dir := filepath.Dir(desktopPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create applications directory: %w", err)
	}

	// Generate .desktop file content
	content := r.generateDesktopFile()

	// Write file
	if err := os.WriteFile(desktopPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write desktop file: %w", err)
	}

	return nil
}

// generateDesktopFile generates the .desktop file content
func (r *linuxRegistrar) generateDesktopFile() string {
	return fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=Switchboard
Comment=Smart URL router for opening links in different browsers
Exec=%s open %%u
Terminal=false
Categories=Network;WebBrowser;
MimeType=x-scheme-handler/http;x-scheme-handler/https;text/html;
StartupNotify=true
`, r.binaryPath)
}

// updateDesktopDatabase updates the desktop database cache
func (r *linuxRegistrar) updateDesktopDatabase() error {
	home, _ := os.UserHomeDir()
	appsDir := filepath.Join(home, ".local", "share", "applications")

	// Check if update-desktop-database is available
	if _, err := exec.LookPath("update-desktop-database"); err != nil {
		// Not available, skip
		return nil
	}

	cmd := exec.Command("update-desktop-database", appsDir)
	return cmd.Run()
}

// setAsDefaultBrowser attempts to set Switchboard as the default browser
func (r *linuxRegistrar) setAsDefaultBrowser() error {
	// Check if xdg-settings is available
	if _, err := exec.LookPath("xdg-settings"); err != nil {
		return err
	}

	// Set as default browser
	cmd := exec.Command("xdg-settings", "set", "default-web-browser", "switchboard.desktop")
	return cmd.Run()
}
