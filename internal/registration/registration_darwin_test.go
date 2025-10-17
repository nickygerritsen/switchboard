//go:build darwin

package registration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewDarwinRegistrar(t *testing.T) {
	reg, err := newDarwinRegistrar()
	if err != nil {
		t.Fatalf("newDarwinRegistrar() returned error: %v", err)
	}

	if reg == nil {
		t.Fatal("newDarwinRegistrar() returned nil")
	}

	if reg.binaryPath == "" {
		t.Error("newDarwinRegistrar() created registrar with empty binaryPath")
	}
}

func TestDarwinRegistrarGetBinaryPath(t *testing.T) {
	reg, err := newDarwinRegistrar()
	if err != nil {
		t.Fatalf("newDarwinRegistrar() returned error: %v", err)
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

func TestDarwinRegistrarGetAppBundlePath(t *testing.T) {
	reg := &darwinRegistrar{binaryPath: "/usr/local/bin/switchboard"}

	appPath := reg.getAppBundlePath()

	if appPath == "" {
		t.Error("getAppBundlePath() returned empty path")
	}

	if !strings.HasSuffix(appPath, "Switchboard.app") {
		t.Errorf("getAppBundlePath() should end with Switchboard.app, got: %s", appPath)
	}

	if !strings.Contains(appPath, "Applications") {
		t.Errorf("getAppBundlePath() should contain Applications, got: %s", appPath)
	}

	home, _ := os.UserHomeDir()
	expectedPath := filepath.Join(home, "Applications", "Switchboard.app")
	if appPath != expectedPath {
		t.Errorf("getAppBundlePath() = %s, want %s", appPath, expectedPath)
	}
}

func TestDarwinRegistrarIsAppBundle(t *testing.T) {
	tests := []struct {
		name       string
		binaryPath string
		want       bool
	}{
		{
			name:       "is app bundle",
			binaryPath: "/Applications/Switchboard.app/Contents/MacOS/Switchboard",
			want:       true,
		},
		{
			name:       "is user app bundle",
			binaryPath: "/Users/test/Applications/Switchboard.app/Contents/MacOS/Switchboard",
			want:       true,
		},
		{
			name:       "not app bundle - regular binary",
			binaryPath: "/usr/local/bin/switchboard",
			want:       false,
		},
		{
			name:       "not app bundle - relative path",
			binaryPath: "./switchboard",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &darwinRegistrar{binaryPath: tt.binaryPath}
			got := reg.isAppBundle()
			if got != tt.want {
				t.Errorf("isAppBundle() = %v, want %v for path %s", got, tt.want, tt.binaryPath)
			}
		})
	}
}

func TestDarwinRegistrarGenerateInfoPlist(t *testing.T) {
	reg := &darwinRegistrar{binaryPath: "/usr/local/bin/switchboard"}

	plist := reg.generateInfoPlist()

	if plist == "" {
		t.Error("generateInfoPlist() returned empty string")
	}

	// Check for required keys
	requiredKeys := []string{
		"CFBundleIdentifier",
		"com.github.nickygerritsen.switchboard",
		"CFBundleName",
		"Switchboard",
		"CFBundleURLTypes",
		"http",
		"https",
		"CFBundleExecutable",
		"LSUIElement",
	}

	for _, key := range requiredKeys {
		if !strings.Contains(plist, key) {
			t.Errorf("generateInfoPlist() missing required key or value: %s", key)
		}
	}

	// Should be valid XML
	if !strings.HasPrefix(plist, "<?xml") {
		t.Error("generateInfoPlist() should start with XML declaration")
	}

	if !strings.Contains(plist, "<!DOCTYPE plist") {
		t.Error("generateInfoPlist() should contain plist DOCTYPE")
	}
}

func TestDarwinRegistrarSwiftSourceGeneration(t *testing.T) {
	testBinaryPath := "/usr/local/bin/test-switchboard"

	// Generate the Swift source that would be used
	swiftSource := generateSwiftSource(testBinaryPath)

	// Check for required Swift components
	requiredComponents := []string{
		"import Cocoa",
		"class URLHandler",
		"NSAppleEventManager",
		"kInternetEventClass",
		"kAEGetURL",
		"NSApplication.shared",
		"setActivationPolicy(.prohibited)",
		testBinaryPath,
		`task.arguments = ["open", urlString]`,
	}

	for _, component := range requiredComponents {
		if !strings.Contains(swiftSource, component) {
			t.Errorf("Swift source missing required component: %s", component)
		}
	}

	// Should not contain any format placeholders
	if strings.Contains(swiftSource, "%s") {
		t.Error("Swift source contains unresolved format placeholders")
	}
}

// generateSwiftSource is extracted for testing purposes
func generateSwiftSource(binaryPath string) string {
	return fmt.Sprintf(`import Cocoa

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
`, binaryPath)
}
