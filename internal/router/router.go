package router

import (
	"github.com/nickygerritsen/switchboard/internal/config"
	"github.com/nickygerritsen/switchboard/internal/logger"
	"github.com/nickygerritsen/switchboard/internal/matcher"
	"github.com/nickygerritsen/switchboard/internal/rewriter"
)

// Router handles URL routing based on configured rules
type Router struct {
	config *config.Config
}

// NewRouter creates a new router instance
func NewRouter(cfg *config.Config) *Router {
	return &Router{
		config: cfg,
	}
}

// FindMatch finds the first matching rule for a URL and returns the browser name, profile, incognito flag, rewritten URL, and whether a match was found
// If no match is found, returns the default browser with empty profile, false for incognito, original URL, and false for matched
func (r *Router) FindMatch(url string) (browser, profile string, incognito, matched bool, rewrittenURL string) {
	logger.Debug("Finding match for URL: %s", url)

	// Try to match against each rule in order (first match wins)
	for i, rule := range r.config.Rules {
		// Check each pattern in the rule
		for _, pattern := range rule.Match {
			if matcher.Match(pattern, url) {
				// Apply URL rewriting if configured
				finalURL := url
				if rule.Rewrite != "" {
					rewritten, err := rewriter.Rewrite(url, rule.Rewrite)
					if err != nil {
						logger.Warn("Failed to rewrite URL %s with template %s: %v", url, rule.Rewrite, err)
					} else {
						finalURL = rewritten
						logger.Info("URL %s rewritten to %s", url, finalURL)
					}
				}

				logger.Info("URL %s matched rule #%d (pattern: %s) -> browser: %s, profile: %s, incognito: %v",
					url, i+1, pattern, rule.Browser, rule.Profile, rule.Incognito)
				return rule.Browser, rule.Profile, rule.Incognito, true, finalURL
			}
		}
	}

	// No match found, use default browser
	logger.Info("URL %s did not match any rules, using default browser: %s", url, r.config.DefaultBrowser)
	return r.config.DefaultBrowser, "", false, false, url
}
