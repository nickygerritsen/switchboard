package browser

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestDetectProfiles_Chromium(t *testing.T) {
	// Create a temporary Chrome-like profile directory structure
	tempDir := t.TempDir()

	// Create "Default" profile
	defaultDir := filepath.Join(tempDir, "Default")
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create Preferences file for Default profile
	prefsDefault := `{"profile": {"name": "Person 1"}}`
	if err := os.WriteFile(filepath.Join(defaultDir, "Preferences"), []byte(prefsDefault), 0644); err != nil {
		t.Fatal(err)
	}

	// Create "Profile 1" profile
	profile1Dir := filepath.Join(tempDir, "Profile 1")
	if err := os.MkdirAll(profile1Dir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create Preferences file for Profile 1
	prefsProfile1 := `{"profile": {"name": "Work"}}`
	if err := os.WriteFile(filepath.Join(profile1Dir, "Preferences"), []byte(prefsProfile1), 0644); err != nil {
		t.Fatal(err)
	}

	// Create "Profile 2" profile
	profile2Dir := filepath.Join(tempDir, "Profile 2")
	if err := os.MkdirAll(profile2Dir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create Preferences file for Profile 2
	prefsProfile2 := `{"profile": {"name": "Personal"}}`
	if err := os.WriteFile(filepath.Join(profile2Dir, "Preferences"), []byte(prefsProfile2), 0644); err != nil {
		t.Fatal(err)
	}

	// Mock getProfileDir to return our temp directory
	originalGetProfileDir := getProfileDirFunc
	defer func() { getProfileDirFunc = originalGetProfileDir }()
	getProfileDirFunc = func(browserName string) string {
		return tempDir
	}

	profiles := detectChromiumProfiles("chrome")

	if len(profiles) != 3 {
		t.Errorf("Expected 3 profiles, got %d", len(profiles))
	}

	// Check Default profile
	if profiles[0].Name != "Person 1" {
		t.Errorf("Expected 'Person 1', got '%s'", profiles[0].Name)
	}
	if profiles[0].Directory != "Default" {
		t.Errorf("Expected 'Default', got '%s'", profiles[0].Directory)
	}

	// Check Profile 1
	if profiles[1].Name != "Work" {
		t.Errorf("Expected 'Work', got '%s'", profiles[1].Name)
	}
	if profiles[1].Directory != "Profile 1" {
		t.Errorf("Expected 'Profile 1', got '%s'", profiles[1].Directory)
	}

	// Check Profile 2
	if profiles[2].Name != "Personal" {
		t.Errorf("Expected 'Personal', got '%s'", profiles[2].Name)
	}
	if profiles[2].Directory != "Profile 2" {
		t.Errorf("Expected 'Profile 2', got '%s'", profiles[2].Directory)
	}
}

func TestDetectProfiles_Chromium_NoPreferences(t *testing.T) {
	// Create a temporary Chrome-like profile directory structure without Preferences
	tempDir := t.TempDir()

	// Create "Default" profile without Preferences file
	defaultDir := filepath.Join(tempDir, "Default")
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Mock getProfileDir to return our temp directory
	originalGetProfileDir := getProfileDirFunc
	defer func() { getProfileDirFunc = originalGetProfileDir }()
	getProfileDirFunc = func(browserName string) string {
		return tempDir
	}

	profiles := detectChromiumProfiles("chrome")

	if len(profiles) != 1 {
		t.Errorf("Expected 1 profile, got %d", len(profiles))
	}

	// Should use directory name as fallback
	if profiles[0].Name != "Default" {
		t.Errorf("Expected 'Default', got '%s'", profiles[0].Name)
	}
}

func TestDetectProfiles_Firefox(t *testing.T) {
	// Create a temporary Firefox-like profile directory structure
	tempDir := t.TempDir()

	// Create profiles.ini file
	profilesIni := `[Profile0]
Name=default-release
IsRelative=1
Path=abc123.default-release
Default=1

[Profile1]
Name=dev
IsRelative=1
Path=xyz456.dev

[General]
StartWithLastProfile=1
Version=2
`
	if err := os.WriteFile(filepath.Join(tempDir, "profiles.ini"), []byte(profilesIni), 0644); err != nil {
		t.Fatal(err)
	}

	// Mock getProfileDir to return our temp directory
	originalGetProfileDir := getProfileDirFunc
	defer func() { getProfileDirFunc = originalGetProfileDir }()
	getProfileDirFunc = func(browserName string) string {
		return tempDir
	}

	profiles := detectFirefoxProfiles()

	if len(profiles) != 2 {
		t.Errorf("Expected 2 profiles, got %d", len(profiles))
	}

	// Check first profile
	if profiles[0].Name != "default-release" {
		t.Errorf("Expected 'default-release', got '%s'", profiles[0].Name)
	}
	if profiles[0].Directory != "abc123.default-release" {
		t.Errorf("Expected 'abc123.default-release', got '%s'", profiles[0].Directory)
	}

	// Check second profile
	if profiles[1].Name != "dev" {
		t.Errorf("Expected 'dev', got '%s'", profiles[1].Name)
	}
	if profiles[1].Directory != "xyz456.dev" {
		t.Errorf("Expected 'xyz456.dev', got '%s'", profiles[1].Directory)
	}
}

func TestParseFirefoxProfilesIni(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []Profile
	}{
		{
			name: "two profiles",
			content: `[Profile0]
Name=default
Path=abc.default

[Profile1]
Name=work
Path=xyz.work
`,
			expected: []Profile{
				{Name: "default", Directory: "abc.default"},
				{Name: "work", Directory: "xyz.work"},
			},
		},
		{
			name: "profile with extra fields",
			content: `[Profile0]
Name=default
IsRelative=1
Path=abc.default
Default=1

[General]
StartWithLastProfile=1
`,
			expected: []Profile{
				{Name: "default", Directory: "abc.default"},
			},
		},
		{
			name:     "empty content",
			content:  "",
			expected: []Profile{},
		},
		{
			name: "profile without name",
			content: `[Profile0]
Path=abc.default
`,
			expected: []Profile{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profiles := parseFirefoxProfilesIni(tt.content)

			if len(profiles) != len(tt.expected) {
				t.Errorf("Expected %d profiles, got %d", len(tt.expected), len(profiles))
				return
			}

			for i, profile := range profiles {
				if profile.Name != tt.expected[i].Name {
					t.Errorf("Profile %d: expected name '%s', got '%s'", i, tt.expected[i].Name, profile.Name)
				}
				if profile.Directory != tt.expected[i].Directory {
					t.Errorf("Profile %d: expected directory '%s', got '%s'", i, tt.expected[i].Directory, profile.Directory)
				}
			}
		})
	}
}

func TestDetectProfiles_UnsupportedBrowser(t *testing.T) {
	profiles := DetectProfiles("safari")
	if profiles != nil {
		t.Errorf("Expected nil profiles for Safari, got %d profiles", len(profiles))
	}

	profiles = DetectProfiles("unknown")
	if profiles != nil {
		t.Errorf("Expected nil profiles for unknown browser, got %d profiles", len(profiles))
	}
}

func TestDetectFirefoxProfiles_WithProfileGroups(t *testing.T) {
	// Simulate a Firefox profile directory with a profiles.ini that
	// references a StoreID, plus a matching Profile Groups SQLite store
	// that names the existing profile and adds another one not in the ini.
	tempDir := t.TempDir()

	profilesIni := `[Profile0]
Name=default-release
IsRelative=1
Path=abc123.default-release
StoreID=deadbeef
ShowSelector=1

[General]
StartWithLastProfile=1
Version=2
`
	if err := os.WriteFile(filepath.Join(tempDir, "profiles.ini"), []byte(profilesIni), 0644); err != nil {
		t.Fatal(err)
	}

	groupDir := filepath.Join(tempDir, "Profile Groups")
	if err := os.MkdirAll(groupDir, 0755); err != nil {
		t.Fatal(err)
	}
	dbPath := filepath.Join(groupDir, "deadbeef.sqlite")
	createFirefoxGroupDB(t, dbPath, []firefoxGroupProfile{
		{Directory: "abc123.default-release", Name: "Privé"},
		{Directory: "xyz999.Profile 1", Name: "Werk"},
	})

	originalGetProfileDir := getProfileDirFunc
	defer func() { getProfileDirFunc = originalGetProfileDir }()
	getProfileDirFunc = func(string) string { return tempDir }

	profiles := detectFirefoxProfiles()
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d: %#v", len(profiles), profiles)
	}

	// Existing ini profile keeps its IniName but takes the SQLite display name.
	if profiles[0].Name != "Privé" {
		t.Errorf("expected display name 'Privé', got %q", profiles[0].Name)
	}
	if profiles[0].IniName != "default-release" {
		t.Errorf("expected IniName 'default-release', got %q", profiles[0].IniName)
	}
	if profiles[0].Path != filepath.Join(tempDir, "abc123.default-release") {
		t.Errorf("unexpected Path: %s", profiles[0].Path)
	}

	// New-style profile (only in SQLite) has no IniName.
	if profiles[1].Name != "Werk" {
		t.Errorf("expected display name 'Werk', got %q", profiles[1].Name)
	}
	if profiles[1].IniName != "" {
		t.Errorf("expected empty IniName for new-style profile, got %q", profiles[1].IniName)
	}
	if profiles[1].Path != filepath.Join(tempDir, "xyz999.Profile 1") {
		t.Errorf("unexpected Path: %s", profiles[1].Path)
	}
}

func TestMergeFirefoxProfiles_NoStoreID(t *testing.T) {
	// Without a StoreID, the SQLite reader must not be consulted at all.
	called := false
	readGroup := func(profileDir, storeID string) []firefoxGroupProfile {
		called = true
		return nil
	}
	profiles := mergeFirefoxProfiles(
		"/firefox",
		[]firefoxIniProfile{
			{Name: "default", Directory: "abc.default"},
		},
		readGroup,
	)
	if called {
		t.Error("readGroup should not be invoked when no profile has a StoreID")
	}
	if len(profiles) != 1 || profiles[0].IniName != "default" || profiles[0].Name != "default" {
		t.Errorf("unexpected merge result: %#v", profiles)
	}
}

// createFirefoxGroupDB writes a minimal Profile Groups SQLite store
// containing the given profile rows, mirroring the schema Firefox uses.
func createFirefoxGroupDB(t *testing.T, path string, rows []firefoxGroupProfile) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer func() { _ = db.Close() }()

	if _, err := db.Exec(`CREATE TABLE Profiles (
		id INTEGER NOT NULL PRIMARY KEY,
		path TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		avatar TEXT NOT NULL,
		themeId TEXT NOT NULL,
		themeFg TEXT NOT NULL,
		themeBg TEXT NOT NULL
	)`); err != nil {
		t.Fatalf("create table: %v", err)
	}

	for i, r := range rows {
		if _, err := db.Exec(
			`INSERT INTO Profiles(id, path, name, avatar, themeId, themeFg, themeBg) VALUES (?, ?, ?, '', '', '', '')`,
			i+1, r.Directory, r.Name,
		); err != nil {
			t.Fatalf("insert row %d: %v", i, err)
		}
	}
}

func TestDetectProfiles_NonexistentDirectory(t *testing.T) {
	// Mock getProfileDir to return a non-existent directory
	originalGetProfileDir := getProfileDirFunc
	defer func() { getProfileDirFunc = originalGetProfileDir }()
	getProfileDirFunc = func(browserName string) string {
		return "/nonexistent/directory/that/does/not/exist"
	}

	profiles := detectChromiumProfiles("chrome")
	if profiles != nil {
		t.Errorf("Expected nil profiles for non-existent directory, got %d profiles", len(profiles))
	}

	profiles = detectFirefoxProfiles()
	if profiles != nil {
		t.Errorf("Expected nil profiles for non-existent directory, got %d profiles", len(profiles))
	}
}
