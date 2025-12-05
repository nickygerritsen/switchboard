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
	// Check if Swift compiler is available
	if _, err := exec.LookPath("swiftc"); err != nil {
		return fmt.Errorf("swift compiler (swiftc) not found, please install Xcode Command Line Tools:\n\n" +
			"    xcode-select --install\n\n" +
			"or install the full Xcode from the App Store\n" +
			"(the Swift compiler is required to build the macOS app bundle for URL handling)")
	}

	fmt.Println("Creating macOS app bundle (this may take a moment to compile)...")

	appPath := r.getAppBundlePath()
	contentsPath := filepath.Join(appPath, "Contents")
	macosPath := filepath.Join(contentsPath, "MacOS")
	resourcesPath := filepath.Join(contentsPath, "Resources")

	// Create directory structure
	if err := os.MkdirAll(macosPath, 0755); err != nil {
		return fmt.Errorf("failed to create app bundle structure: %w", err)
	}
	if err := os.MkdirAll(resourcesPath, 0755); err != nil {
		return fmt.Errorf("failed to create resources directory: %w", err)
	}

	// Create Info.plist
	infoPlist := r.generateInfoPlist()
	plistPath := filepath.Join(contentsPath, "Info.plist")
	if err := os.WriteFile(plistPath, []byte(infoPlist), 0644); err != nil {
		return fmt.Errorf("failed to write Info.plist: %w", err)
	}

	// Create native Swift launcher that handles Apple Events
	// This is fast, modern, and doesn't show any GUI
	swiftSource := fmt.Sprintf(`import Cocoa

class URLHandler: NSObject {
    let binaryPath: String

    init(binaryPath: String) {
        self.binaryPath = binaryPath
    }

    @objc func handleGetURL(_ event: NSAppleEventDescriptor, replyEvent: NSAppleEventDescriptor) {
        guard let urlString = event.paramDescriptor(forKeyword: keyDirectObject)?.stringValue else {
            NSApplication.shared.terminate(nil)
            return
        }

        let task = Process()
        task.executableURL = URL(fileURLWithPath: binaryPath)
        task.arguments = ["open", urlString]

        try? task.run()

        // Exit immediately after launching
        NSApplication.shared.terminate(nil)
    }
}

// Initialize the application
let app = NSApplication.shared
app.setActivationPolicy(.prohibited)

let handler = URLHandler(binaryPath: "%s")

NSAppleEventManager.shared().setEventHandler(
    handler,
    andSelector: #selector(URLHandler.handleGetURL(_:replyEvent:)),
    forEventClass: AEEventClass(kInternetEventClass),
    andEventID: AEEventID(kAEGetURL)
)

// Run the event loop
app.run()
`, r.binaryPath)

	// Write Swift source
	sourcePath := filepath.Join(resourcesPath, "launcher.swift")
	if err := os.WriteFile(sourcePath, []byte(swiftSource), 0644); err != nil {
		return fmt.Errorf("failed to write Swift source: %w", err)
	}

	// Compile Swift to native executable
	launcherPath := filepath.Join(macosPath, "Switchboard")
	cmd := exec.Command("swiftc", "-o", launcherPath, sourcePath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to compile Swift launcher: %w (output: %s)", err, string(output))
	}

	// Ensure executable
	if err := os.Chmod(launcherPath, 0755); err != nil {
		return fmt.Errorf("failed to make launcher executable: %w", err)
	}

	fmt.Println("Registering with macOS Launch Services...")

	// Register with LaunchServices
	if err := r.registerAppBundle(); err != nil {
		return err
	}

	fmt.Println("✓ Successfully created and registered Switchboard.app")
	return nil
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
	<key>LSUIElement</key>
	<true/>
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
