package rewriter

import (
	"net/url"
	"strings"
)

// Rewrite rewrites a URL using a template string with substitution variables.
// Supported variables:
//   - {scheme}   - URL scheme (http, https, etc.)
//   - {host}     - Hostname (e.g., "example.com")
//   - {port}     - Port number (empty if not specified)
//   - {path}     - Path portion (e.g., "/foo/bar")
//   - {query}    - Query string (e.g., "key=value&key2=value2")
//   - {fragment} - Fragment/hash (e.g., "section")
//
// Example:
//
//	Rewrite("https://twitter.com/user/status/123", "xcancel.com{path}")
//	Returns: "https://xcancel.com/user/status/123"
func Rewrite(originalURL, template string) (string, error) {
	if template == "" {
		return originalURL, nil
	}

	// Ensure URL has a scheme for proper parsing
	urlToParse := originalURL
	if !strings.Contains(originalURL, "://") {
		urlToParse = "https://" + originalURL
	}

	// Parse the original URL
	parsed, err := url.Parse(urlToParse)
	if err != nil {
		return originalURL, err
	}

	// Extract URL components
	scheme := parsed.Scheme
	if scheme == "" {
		scheme = "https" // Default to https if no scheme
	}
	host := parsed.Hostname()
	port := parsed.Port()
	path := parsed.Path
	query := parsed.RawQuery
	fragment := parsed.Fragment

	// Build the new URL from the template
	result := template

	// Replace template variables
	result = strings.ReplaceAll(result, "{scheme}", scheme)
	result = strings.ReplaceAll(result, "{host}", host)
	result = strings.ReplaceAll(result, "{port}", port)
	result = strings.ReplaceAll(result, "{path}", path)
	result = strings.ReplaceAll(result, "{query}", query)
	result = strings.ReplaceAll(result, "{fragment}", fragment)

	// If the result doesn't start with a scheme, prepend the original scheme
	if !strings.Contains(result, "://") {
		result = scheme + "://" + result
	}

	// Parse the result to ensure it's a valid URL and normalize it
	finalURL, err := url.Parse(result)
	if err != nil {
		return originalURL, err
	}

	return finalURL.String(), nil
}
