package matcher

import (
	"fmt"
	"net/url"
	"strings"
)

// Match checks if a URL matches a given pattern.
// Pattern format supports:
// - Simple domain: "google.com" (matches any protocol, port, path)
// - With protocol: "https://google.com" (matches only HTTPS)
// - With port: "localhost:3000" (matches only that port)
// - With path: "github.com/user" (path prefix match)
// - Wildcard subdomain: "*.google.com" (matches mail.google.com, drive.google.com, etc)
// - Wildcard: "*" (matches everything)
//
// Examples:
//   - "google.com" matches "https://google.com/search"
//   - "https://google.com" matches only HTTPS, not HTTP
//   - "localhost:3000" matches "http://localhost:3000" but not "http://localhost:8080"
//   - "github.com/user" matches "https://github.com/user/repo" but not "https://github.com/settings"
//   - "*.google.com" matches "https://mail.google.com" but not "https://google.com"
func Match(pattern, urlStr string) bool {
	// Handle empty pattern
	if pattern == "" {
		return false
	}

	// Handle wildcard pattern
	if pattern == "*" {
		return true
	}

	// Parse pattern
	patScheme, patHost, patPort, patPath, err := parsePattern(pattern)
	if err != nil {
		return false
	}

	// Parse URL
	u, err := url.Parse(urlStr)
	if err != nil || u.Host == "" {
		return false
	}

	// Extract URL components
	urlScheme := strings.ToLower(u.Scheme)
	urlHost := strings.ToLower(u.Hostname())
	urlPort := u.Port()
	urlPath := u.Path

	// Match scheme (if pattern specifies one)
	if patScheme != "" && patScheme != urlScheme {
		return false
	}

	// Match port (if pattern specifies one)
	if patPort != "" && patPort != urlPort {
		return false
	}

	// Match host
	if !matchHost(patHost, urlHost) {
		return false
	}

	// Match path (if pattern specifies one)
	if patPath != "" {
		if !matchPath(patPath, urlPath) {
			return false
		}
	}

	return true
}

// parsePattern parses a pattern into its components: scheme, host, port, path
// Returns: scheme, host, port, path, error
func parsePattern(pattern string) (string, string, string, string, error) {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return "", "", "", "", fmt.Errorf("empty pattern")
	}

	// Handle single wildcard
	if pattern == "*" {
		return "", "*", "", "", nil
	}

	var scheme, host, port, path string

	// Check if pattern has a scheme
	if strings.Contains(pattern, "://") {
		parts := strings.SplitN(pattern, "://", 2)
		scheme = strings.ToLower(parts[0])
		pattern = parts[1]
	}

	// Check if pattern has a path
	if strings.Contains(pattern, "/") {
		parts := strings.SplitN(pattern, "/", 2)
		host = parts[0]
		path = "/" + parts[1]
	} else {
		host = pattern
	}

	// Check if host has a port
	if strings.Contains(host, ":") && !strings.HasPrefix(host, "[") {
		// Handle wildcard subdomain patterns like "*.example.com:3000"
		// Split from the right to avoid issues with "*." prefix
		lastColon := strings.LastIndex(host, ":")
		port = host[lastColon+1:]
		host = host[:lastColon]
	}

	host = strings.ToLower(host)

	return scheme, host, port, path, nil
}

// matchHost checks if a hostname matches a pattern with wildcard support
func matchHost(pattern, hostname string) bool {
	if pattern == "" || hostname == "" {
		return false
	}

	// Wildcard matches everything
	if pattern == "*" {
		return true
	}

	pattern = strings.ToLower(pattern)
	hostname = strings.ToLower(hostname)

	// Exact match
	if pattern == hostname {
		return true
	}

	// Wildcard subdomain matching (e.g., *.google.com)
	if strings.HasPrefix(pattern, "*.") {
		baseDomain := pattern[2:] // Remove "*."
		// Hostname must end with the base domain and have at least one subdomain
		// e.g., "mail.google.com" matches "*.google.com" but "google.com" does not
		if strings.HasSuffix(hostname, "."+baseDomain) {
			return true
		}
	}

	return false
}

// matchPath checks if a URL path matches a pattern path (prefix matching)
func matchPath(patternPath, urlPath string) bool {
	// Empty pattern path matches any URL path
	if patternPath == "" {
		return true
	}

	// Ensure both paths start with /
	if !strings.HasPrefix(patternPath, "/") {
		patternPath = "/" + patternPath
	}
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}

	// Handle wildcard path (e.g., "/api/*")
	if strings.HasSuffix(patternPath, "/*") {
		prefix := patternPath[:len(patternPath)-1] // Remove the "*", keep the "/"
		return strings.HasPrefix(urlPath, prefix)
	}

	// Handle trailing slash in pattern
	if patternPath == "/" {
		return true // Root path matches everything
	}

	// Remove trailing slash from pattern for comparison (unless it's just "/")
	patternPath = strings.TrimSuffix(patternPath, "/")

	// Prefix match (case-sensitive for paths)
	// Pattern "/user" matches "/user", "/user/", "/user/repo", etc.
	if urlPath == patternPath || strings.HasPrefix(urlPath, patternPath+"/") {
		return true
	}

	return false
}
