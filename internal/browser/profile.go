package browser

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "modernc.org/sqlite"
)

// Profile represents a browser profile
type Profile struct {
	// Name is the display name shown to the user. For Firefox new-style
	// profiles this comes from the Profile Groups SQLite store; otherwise
	// it is the name from profiles.ini (or Preferences for Chromium).
	Name string
	// Directory is the profile's directory name (relative to the browser
	// profile root).
	Directory string
	// Path is the absolute path to the profile directory.
	Path string
	// IniName is the profile's Name= field in Firefox's profiles.ini.
	// Empty for new-style Firefox profiles that exist only in a Profile
	// Groups SQLite store (and for Chromium browsers).
	IniName string
}

// getProfileDirFunc is a function variable that can be mocked in tests
var getProfileDirFunc = getProfileDirImpl

// getProfileDir returns the profile directory for a browser
func getProfileDir(browserName string) string {
	return getProfileDirFunc(browserName)
}

// getProfileDirImpl is the actual implementation
func getProfileDirImpl(browserName string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch runtime.GOOS {
	case "darwin":
		return getDarwinProfileDir(homeDir, browserName)
	case "linux":
		return getLinuxProfileDir(homeDir, browserName)
	case "windows":
		return getWindowsProfileDir(browserName)
	default:
		return ""
	}
}

func getDarwinProfileDir(homeDir, browserName string) string {
	appSupportDir := filepath.Join(homeDir, "Library", "Application Support")

	switch browserName {
	case "chrome":
		return filepath.Join(appSupportDir, "Google", "Chrome")
	case "brave":
		return filepath.Join(appSupportDir, "BraveSoftware", "Brave-Browser")
	case "edge":
		return filepath.Join(appSupportDir, "Microsoft Edge")
	case "chromium":
		return filepath.Join(appSupportDir, "Chromium")
	case "vivaldi":
		return filepath.Join(appSupportDir, "Vivaldi")
	case "firefox":
		return filepath.Join(appSupportDir, "Firefox")
	default:
		return ""
	}
}

func getLinuxProfileDir(homeDir, browserName string) string {
	configDir := filepath.Join(homeDir, ".config")

	switch browserName {
	case "chrome":
		return filepath.Join(configDir, "google-chrome")
	case "brave":
		// Try different possible locations
		dirs := []string{
			filepath.Join(configDir, "BraveSoftware", "Brave-Browser"),
			filepath.Join(configDir, "brave"),
		}
		for _, dir := range dirs {
			if _, err := os.Stat(dir); err == nil {
				return dir
			}
		}
		return ""
	case "edge":
		return filepath.Join(configDir, "microsoft-edge")
	case "chromium":
		return filepath.Join(configDir, "chromium")
	case "vivaldi":
		return filepath.Join(configDir, "vivaldi")
	case "firefox":
		return filepath.Join(homeDir, ".mozilla", "firefox")
	default:
		return ""
	}
}

func getWindowsProfileDir(browserName string) string {
	localAppData := os.Getenv("LOCALAPPDATA")
	appData := os.Getenv("APPDATA")

	if localAppData == "" && appData == "" {
		return ""
	}

	switch browserName {
	case "chrome":
		return filepath.Join(localAppData, "Google", "Chrome", "User Data")
	case "brave":
		return filepath.Join(localAppData, "BraveSoftware", "Brave-Browser", "User Data")
	case "edge":
		return filepath.Join(localAppData, "Microsoft", "Edge", "User Data")
	case "chromium":
		return filepath.Join(localAppData, "Chromium", "User Data")
	case "vivaldi":
		return filepath.Join(localAppData, "Vivaldi", "User Data")
	case "firefox":
		return filepath.Join(appData, "Mozilla", "Firefox", "Profiles")
	default:
		return ""
	}
}

// DetectProfiles detects profiles for a browser
func DetectProfiles(browserName string) []Profile {
	switch browserName {
	case "chrome", "brave", "edge", "chromium", "vivaldi":
		return detectChromiumProfiles(browserName)
	case "firefox":
		return detectFirefoxProfiles()
	default:
		return nil
	}
}

// detectChromiumProfiles detects profiles for Chromium-based browsers
func detectChromiumProfiles(browserName string) []Profile {
	profileDir := getProfileDir(browserName)
	if profileDir == "" {
		return nil
	}

	// Check if profile directory exists
	if _, err := os.Stat(profileDir); os.IsNotExist(err) {
		return nil
	}

	var profiles []Profile

	// Check for "Default" profile
	defaultProfile := filepath.Join(profileDir, "Default")
	if _, err := os.Stat(defaultProfile); err == nil {
		name := getChromiumProfileName(defaultProfile)
		if name == "" {
			name = "Default"
		}
		profiles = append(profiles, Profile{
			Name:      name,
			Directory: "Default",
			Path:      defaultProfile,
		})
	}

	// Check for additional profiles (Profile 1, Profile 2, etc.)
	entries, err := os.ReadDir(profileDir)
	if err != nil {
		return profiles
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Check if it's a profile directory (Profile 1, Profile 2, etc.)
		if strings.HasPrefix(name, "Profile ") {
			profilePath := filepath.Join(profileDir, name)
			profileName := getChromiumProfileName(profilePath)
			if profileName == "" {
				profileName = name
			}
			profiles = append(profiles, Profile{
				Name:      profileName,
				Directory: name,
				Path:      profilePath,
			})
		}
	}

	return profiles
}

// getChromiumProfileName reads the profile name from Preferences file
func getChromiumProfileName(profilePath string) string {
	prefsPath := filepath.Join(profilePath, "Preferences")
	data, err := os.ReadFile(prefsPath)
	if err != nil {
		return ""
	}

	var prefs struct {
		Profile struct {
			Name string `json:"name"`
		} `json:"profile"`
	}

	if err := json.Unmarshal(data, &prefs); err != nil {
		return ""
	}

	return prefs.Profile.Name
}

// firefoxIniProfile is an entry parsed out of Firefox's profiles.ini.
type firefoxIniProfile struct {
	Name      string
	Directory string
	StoreID   string
}

// firefoxGroupProfile is an entry parsed out of a Profile Groups SQLite
// store. Each row represents a profile that belongs to the same Firefox
// installation group.
type firefoxGroupProfile struct {
	Directory string
	Name      string
}

// detectFirefoxProfiles detects Firefox profiles, including new-style
// profiles tracked in the Profile Groups SQLite store (Firefox 138+).
func detectFirefoxProfiles() []Profile {
	profileDir := getProfileDir("firefox")
	if profileDir == "" {
		return nil
	}

	// Check if profile directory exists
	if _, err := os.Stat(profileDir); os.IsNotExist(err) {
		return nil
	}

	profilesIni := filepath.Join(profileDir, "profiles.ini")
	data, err := os.ReadFile(profilesIni)
	if err != nil {
		return nil
	}

	iniProfiles := parseFirefoxProfilesIni(string(data))
	return mergeFirefoxProfiles(profileDir, iniProfiles, readFirefoxGroupProfiles)
}

// mergeFirefoxProfiles combines profiles from profiles.ini with any extra
// profiles tracked in the Profile Groups SQLite stores referenced by
// StoreID. SQLite names take precedence for display.
func mergeFirefoxProfiles(
	profileDir string,
	iniProfiles []firefoxIniProfile,
	readGroup func(profileDir, storeID string) []firefoxGroupProfile,
) []Profile {
	// Index ini profiles by directory for quick lookups when merging.
	type entry struct {
		profile Profile
		fromIni bool
	}
	byDir := make(map[string]*entry)
	var order []string

	for _, ip := range iniProfiles {
		if ip.Directory == "" {
			continue
		}
		p := Profile{
			Name:      ip.Name,
			Directory: ip.Directory,
			Path:      filepath.Join(profileDir, ip.Directory),
			IniName:   ip.Name,
		}
		byDir[ip.Directory] = &entry{profile: p, fromIni: true}
		order = append(order, ip.Directory)
	}

	// Collect unique StoreIDs referenced from the ini.
	seenStore := make(map[string]bool)
	for _, ip := range iniProfiles {
		if ip.StoreID == "" || seenStore[ip.StoreID] {
			continue
		}
		seenStore[ip.StoreID] = true

		for _, gp := range readGroup(profileDir, ip.StoreID) {
			if gp.Directory == "" {
				continue
			}
			if existing, ok := byDir[gp.Directory]; ok {
				// Prefer the SQLite display name when available.
				if gp.Name != "" {
					existing.profile.Name = gp.Name
				}
				continue
			}
			byDir[gp.Directory] = &entry{
				profile: Profile{
					Name:      gp.Name,
					Directory: gp.Directory,
					Path:      filepath.Join(profileDir, gp.Directory),
				},
			}
			order = append(order, gp.Directory)
		}
	}

	var profiles []Profile
	for _, dir := range order {
		e := byDir[dir]
		// Skip entries whose ini section had no Name= (preserves the
		// behaviour of the previous implementation, which required Name).
		if e.fromIni && e.profile.IniName == "" {
			continue
		}
		if !e.fromIni && e.profile.Name == "" {
			continue
		}
		profiles = append(profiles, e.profile)
	}
	return profiles
}

// parseFirefoxProfilesIni parses Firefox's profiles.ini file.
func parseFirefoxProfilesIni(content string) []firefoxIniProfile {
	var profiles []firefoxIniProfile
	lines := strings.Split(content, "\n")

	var currentProfile *firefoxIniProfile
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Start of a profile section
		if strings.HasPrefix(line, "[Profile") {
			if currentProfile != nil {
				profiles = append(profiles, *currentProfile)
			}
			currentProfile = &firefoxIniProfile{}
			continue
		}

		// End of profiles.ini or start of different section
		if strings.HasPrefix(line, "[") && !strings.HasPrefix(line, "[Profile") {
			if currentProfile != nil {
				profiles = append(profiles, *currentProfile)
				currentProfile = nil
			}
			continue
		}

		if currentProfile == nil {
			continue
		}

		// Parse key=value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Name":
			currentProfile.Name = value
		case "Path":
			currentProfile.Directory = value
		case "StoreID":
			currentProfile.StoreID = value
		}
	}

	// Add the last profile
	if currentProfile != nil && currentProfile.Name != "" {
		profiles = append(profiles, *currentProfile)
	}

	return profiles
}

// readFirefoxGroupProfiles reads the Profiles table from a Profile Groups
// SQLite store for the given StoreID. Returns an empty slice if the file
// doesn't exist or can't be read — this is best-effort enrichment.
func readFirefoxGroupProfiles(profileDir, storeID string) []firefoxGroupProfile {
	dbPath := filepath.Join(profileDir, "Profile Groups", storeID+".sqlite")
	if _, err := os.Stat(dbPath); err != nil {
		return nil
	}

	// Open read-only so we don't compete with a running Firefox for write
	// locks. modernc.org/sqlite honours the standard URI query options.
	db, err := sql.Open("sqlite", "file:"+dbPath+"?mode=ro&immutable=1")
	if err != nil {
		return nil
	}
	defer func() { _ = db.Close() }()

	rows, err := db.Query("SELECT path, name FROM Profiles")
	if err != nil {
		return nil
	}
	defer func() { _ = rows.Close() }()

	var result []firefoxGroupProfile
	for rows.Next() {
		var path, name string
		if err := rows.Scan(&path, &name); err != nil {
			continue
		}
		result = append(result, firefoxGroupProfile{Directory: path, Name: name})
	}
	return result
}
