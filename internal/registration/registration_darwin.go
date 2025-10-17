//go:build darwin

package registration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// New creates a new Registrar for the current platform
func New() (Registrar, error) {
	return newDarwinRegistrar()
}

type darwinRegistrar struct {
	binaryPath string
}

func newDarwinRegistrar() (*darwinRegistrar, error) {
	binPath, err := getBinaryPath()
	if err != nil {
		return nil, err
	}

	return &darwinRegistrar{
		binaryPath: binPath,
	}, nil
}

func (r *darwinRegistrar) Register() error {
	// Check if we're running from an .app bundle
	if r.isAppBundle() {
		return r.registerAppBundle()
	}

	// For CLI binary, we need to create a helper .app bundle
	return r.createAndRegisterAppBundle()
}

func (r *darwinRegistrar) Unregister() error {
	appPath := r.getAppBundlePath()

	// Remove the .app bundle if it exists
	if _, err := os.Stat(appPath); err == nil {
		if err := os.RemoveAll(appPath); err != nil {
			return fmt.Errorf("failed to remove app bundle: %w", err)
		}
	}

	// Rebuild LaunchServices database
	return r.rebuildLaunchServices()
}

func (r *darwinRegistrar) IsRegistered() (bool, error) {
	appPath := r.getAppBundlePath()

	// Check if .app bundle exists
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

func (r *darwinRegistrar) GetBinaryPath() (string, error) {
	return r.binaryPath, nil
}

// isAppBundle checks if the binary is running from a .app bundle
func (r *darwinRegistrar) isAppBundle() bool {
	return strings.Contains(r.binaryPath, ".app/Contents/MacOS/")
}

// getAppBundlePath returns the path where the .app bundle should be
func (r *darwinRegistrar) getAppBundlePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Applications", "Switchboard.app")
}

// createAndRegisterAppBundle creates a minimal .app bundle and registers it
func (r *darwinRegistrar) createAndRegisterAppBundle() error {
	appPath := r.getAppBundlePath()
	contentsPath := filepath.Join(appPath, "Contents")
	macosPath := filepath.Join(contentsPath, "MacOS")

	// Create directory structure
	if err := os.MkdirAll(macosPath, 0755); err != nil {
		return fmt.Errorf("failed to create app bundle structure: %w", err)
	}

	// Create Info.plist
	infoPlist := r.generateInfoPlist()
	plistPath := filepath.Join(contentsPath, "Info.plist")
	if err := os.WriteFile(plistPath, []byte(infoPlist), 0644); err != nil {
		return fmt.Errorf("failed to write Info.plist: %w", err)
	}

	// Create launcher script that calls the actual binary
	launcherScript := fmt.Sprintf(`#!/bin/bash
# Switchboard launcher script
exec "%s" open "$@"
`, r.binaryPath)

	launcherPath := filepath.Join(macosPath, "Switchboard")
	if err := os.WriteFile(launcherPath, []byte(launcherScript), 0755); err != nil {
		return fmt.Errorf("failed to write launcher script: %w", err)
	}

	// Register with LaunchServices
	return r.registerAppBundle()
}

// registerAppBundle registers the .app bundle with LaunchServices
func (r *darwinRegistrar) registerAppBundle() error {
	appPath := r.getAppBundlePath()

	// Use lsregister to register the app
	lsregisterPath := "/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister"
	cmd := exec.Command(lsregisterPath, "-f", appPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to register with LaunchServices: %w", err)
	}

	return nil
}

// rebuildLaunchServices rebuilds the LaunchServices database
func (r *darwinRegistrar) rebuildLaunchServices() error {
	// On modern macOS, simply removing the app bundle is sufficient
	// The LaunchServices database will update automatically
	// The -kill option was removed from lsregister as it's no longer needed
	return nil
}

// generateInfoPlist generates the Info.plist content for the .app bundle
func (r *darwinRegistrar) generateInfoPlist() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleIdentifier</key>
	<string>com.github.nickygerritsen.switchboard</string>
	<key>CFBundleName</key>
	<string>Switchboard</string>
	<key>CFBundleDisplayName</key>
	<string>Switchboard</string>
	<key>CFBundleVersion</key>
	<string>1.0</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleExecutable</key>
	<string>Switchboard</string>
	<key>CFBundleURLTypes</key>
	<array>
		<dict>
			<key>CFBundleURLName</key>
			<string>Web site URL</string>
			<key>CFBundleURLSchemes</key>
			<array>
				<string>http</string>
				<string>https</string>
			</array>
		</dict>
	</array>
	<key>CFBundleDocumentTypes</key>
	<array>
		<dict>
			<key>CFBundleTypeRole</key>
			<string>Viewer</string>
			<key>LSItemContentTypes</key>
			<array>
				<string>public.html</string>
				<string>public.xhtml</string>
			</array>
		</dict>
	</array>
	<key>NSPrincipalClass</key>
	<string>NSApplication</string>
</dict>
</plist>
`
}
