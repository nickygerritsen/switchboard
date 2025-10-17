package main

import (
	"fmt"

	"github.com/nickygerritsen/switchboard/internal/browser"
)

// fakeDetector is a fake browser detector for testing
type fakeDetector struct {
	browsers    map[string]*browser.Browser
	detectError error
}

func (f *fakeDetector) Detect(name string) (*browser.Browser, error) {
	if f.detectError != nil {
		return nil, f.detectError
	}
	if br, ok := f.browsers[name]; ok {
		return br, nil
	}
	return nil, fmt.Errorf("browser not found: %s", name)
}

func (f *fakeDetector) DetectAll() map[string]*browser.Browser {
	return f.browsers
}

// fakeRouter is a fake URL router for testing
type fakeRouter struct {
	matches map[string]routeResult
}

type routeResult struct {
	browser   string
	profile   string
	incognito bool
	matched   bool
}

func (f *fakeRouter) FindMatch(url string) (browser, profile string, incognito, matched bool) {
	if result, ok := f.matches[url]; ok {
		return result.browser, result.profile, result.incognito, result.matched
	}
	// Default behavior
	return "firefox", "", false, false
}

// fakeLauncher is a fake browser launcher for testing
type fakeLauncher struct {
	launchError  error
	launchedURLs []launchRecord
}

type launchRecord struct {
	browser   string
	url       string
	profile   string
	incognito bool
}

func (f *fakeLauncher) Launch(br *browser.Browser, url, profile string, incognito bool) error {
	if f.launchError != nil {
		return f.launchError
	}
	f.launchedURLs = append(f.launchedURLs, launchRecord{
		browser:   br.Name,
		url:       url,
		profile:   profile,
		incognito: incognito,
	})
	return nil
}
