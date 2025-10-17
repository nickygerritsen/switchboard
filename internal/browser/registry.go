package browser

import "runtime"

// Browser represents a detected browser
type Browser struct {
	Name     string
	Path     string
	Profiles []Profile
}

// BrowserDef defines a browser and its potential paths
type BrowserDef struct {
	Name         string
	Aliases      []string
	DarwinPaths  []string
	LinuxPaths   []string
	WindowsPaths []string
}

// knownBrowsers contains the registry of known browsers and their common paths
var knownBrowsers = []BrowserDef{
	{
		Name:    "chrome",
		Aliases: []string{"google-chrome", "google chrome"},
		DarwinPaths: []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		},
		LinuxPaths: []string{
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/snap/bin/chromium",
			"/var/lib/flatpak/exports/bin/com.google.Chrome",
		},
		WindowsPaths: []string{
			"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
			"C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe",
		},
	},
	{
		Name:    "firefox",
		Aliases: []string{"mozilla-firefox", "mozilla firefox"},
		DarwinPaths: []string{
			"/Applications/Firefox.app/Contents/MacOS/firefox",
		},
		LinuxPaths: []string{
			"/usr/bin/firefox",
			"/usr/bin/firefox-esr",
			"/snap/bin/firefox",
			"/var/lib/flatpak/exports/bin/org.mozilla.firefox",
		},
		WindowsPaths: []string{
			"C:\\Program Files\\Mozilla Firefox\\firefox.exe",
			"C:\\Program Files (x86)\\Mozilla Firefox\\firefox.exe",
		},
	},
	{
		Name:    "brave",
		Aliases: []string{"brave-browser", "brave browser"},
		DarwinPaths: []string{
			"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
		},
		LinuxPaths: []string{
			"/usr/bin/brave",
			"/usr/bin/brave-browser",
			"/snap/bin/brave",
			"/var/lib/flatpak/exports/bin/com.brave.Browser",
		},
		WindowsPaths: []string{
			"C:\\Program Files\\BraveSoftware\\Brave-Browser\\Application\\brave.exe",
			"C:\\Program Files (x86)\\BraveSoftware\\Brave-Browser\\Application\\brave.exe",
		},
	},
	{
		Name:    "vivaldi",
		Aliases: []string{},
		DarwinPaths: []string{
			"/Applications/Vivaldi.app/Contents/MacOS/Vivaldi",
		},
		LinuxPaths: []string{
			"/usr/bin/vivaldi",
			"/usr/bin/vivaldi-stable",
			"/snap/bin/vivaldi",
			"/var/lib/flatpak/exports/bin/com.vivaldi.Vivaldi",
		},
		WindowsPaths: []string{
			"C:\\Program Files\\Vivaldi\\Application\\vivaldi.exe",
			"C:\\Program Files (x86)\\Vivaldi\\Application\\vivaldi.exe",
		},
	},
	{
		Name:    "safari",
		Aliases: []string{},
		DarwinPaths: []string{
			"/Applications/Safari.app/Contents/MacOS/Safari",
		},
		LinuxPaths:   []string{}, // Safari not available on Linux
		WindowsPaths: []string{}, // Safari not available on Windows anymore
	},
	{
		Name:    "edge",
		Aliases: []string{"microsoft-edge", "microsoft edge", "msedge"},
		DarwinPaths: []string{
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		},
		LinuxPaths: []string{
			"/usr/bin/microsoft-edge",
			"/usr/bin/microsoft-edge-stable",
			"/usr/bin/microsoft-edge-dev",
			"/var/lib/flatpak/exports/bin/com.microsoft.Edge",
		},
		WindowsPaths: []string{
			"C:\\Program Files\\Microsoft\\Edge\\Application\\msedge.exe",
			"C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
		},
	},
	{
		Name:    "chromium",
		Aliases: []string{},
		DarwinPaths: []string{
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		},
		LinuxPaths: []string{
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/snap/bin/chromium",
			"/var/lib/flatpak/exports/bin/org.chromium.Chromium",
		},
		WindowsPaths: []string{
			"C:\\Program Files\\Chromium\\Application\\chrome.exe",
			"C:\\Program Files (x86)\\Chromium\\Application\\chrome.exe",
		},
	},
}

// GetPaths returns the potential paths for a browser on the current OS
func (b *BrowserDef) GetPaths() []string {
	switch runtime.GOOS {
	case "darwin":
		return b.DarwinPaths
	case "linux":
		return b.LinuxPaths
	case "windows":
		return b.WindowsPaths
	default:
		return []string{}
	}
}

// GetBrowserDef returns the browser definition for a given browser name
func GetBrowserDef(name string) *BrowserDef {
	// Normalize name to lowercase for comparison
	nameLower := name

	for i := range knownBrowsers {
		browser := &knownBrowsers[i]
		if browser.Name == nameLower {
			return browser
		}
		// Check aliases
		for _, alias := range browser.Aliases {
			if alias == nameLower {
				return browser
			}
		}
	}
	return nil
}

// GetAllBrowserDefs returns all known browser definitions
func GetAllBrowserDefs() []BrowserDef {
	return knownBrowsers
}
