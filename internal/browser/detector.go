package browser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/nickygerritsen/switchboard/internal/config"
)

// Detector handles browser detection and caching
type Detector struct {
	config   *config.Config
	cache    map[string]*Browser
	cacheMux sync.RWMutex
}

// NewDetector creates a new browser detector
func NewDetector(cfg *config.Config) *Detector {
	return &Detector{
		config: cfg,
		cache:  make(map[string]*Browser),
	}
}

// Detect finds an installed browser by name
func (d *Detector) Detect(name string) (*Browser, error) {
	// Check cache first
	d.cacheMux.RLock()
	if browser, ok := d.cache[name]; ok {
		d.cacheMux.RUnlock()
		return browser, nil
	}
	d.cacheMux.RUnlock()

	// Check if user specified a custom path in config
	if d.config != nil && d.config.Browsers != nil {
		if browserCfg, ok := d.config.Browsers[name]; ok {
			if browserCfg.Path != "" && browserCfg.Path != "auto" {
				// Verify the path exists
				if _, err := os.Stat(browserCfg.Path); err == nil {
					browser := &Browser{
						Name:     name,
						Path:     browserCfg.Path,
						Profiles: DetectProfiles(name),
					}
					d.cacheResult(name, browser)
					return browser, nil
				}
			}
		}
	}

	// Try to auto-detect
	browserDef := GetBrowserDef(name)
	if browserDef == nil {
		return nil, fmt.Errorf("unknown browser: %s", name)
	}

	// Check known paths for this OS
	for _, path := range browserDef.GetPaths() {
		if _, err := os.Stat(path); err == nil {
			browser := &Browser{
				Name:     name,
				Path:     path,
				Profiles: DetectProfiles(name),
			}
			d.cacheResult(name, browser)
			return browser, nil
		}
	}

	// Try to find in PATH
	path, err := exec.LookPath(name)
	if err == nil {
		absPath, err := filepath.Abs(path)
		if err == nil {
			browser := &Browser{
				Name:     name,
				Path:     absPath,
				Profiles: DetectProfiles(name),
			}
			d.cacheResult(name, browser)
			return browser, nil
		}
	}

	// Try aliases in PATH
	if browserDef.Aliases != nil {
		for _, alias := range browserDef.Aliases {
			path, err := exec.LookPath(alias)
			if err == nil {
				absPath, err := filepath.Abs(path)
				if err == nil {
					browser := &Browser{
						Name:     name,
						Path:     absPath,
						Profiles: DetectProfiles(name),
					}
					d.cacheResult(name, browser)
					return browser, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("browser not found: %s", name)
}

// DetectAll finds all installed browsers
func (d *Detector) DetectAll() map[string]*Browser {
	result := make(map[string]*Browser)

	for _, browserDef := range GetAllBrowserDefs() {
		browser, err := d.Detect(browserDef.Name)
		if err == nil {
			result[browserDef.Name] = browser
		}
	}

	return result
}

// cacheResult stores a detection result in the cache
func (d *Detector) cacheResult(name string, browser *Browser) {
	d.cacheMux.Lock()
	defer d.cacheMux.Unlock()
	d.cache[name] = browser
}

// ClearCache clears the detection cache
func (d *Detector) ClearCache() {
	d.cacheMux.Lock()
	defer d.cacheMux.Unlock()
	d.cache = make(map[string]*Browser)
}

// GetCached returns a cached browser if it exists
func (d *Detector) GetCached(name string) (*Browser, bool) {
	d.cacheMux.RLock()
	defer d.cacheMux.RUnlock()
	browser, ok := d.cache[name]
	return browser, ok
}
