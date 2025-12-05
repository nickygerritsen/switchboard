package matcher

import (
	"testing"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		url     string
		want    bool
	}{
		// Exact domain matches
		{
			name:    "exact domain match",
			pattern: "google.com",
			url:     "https://google.com",
			want:    true,
		},
		{
			name:    "exact domain match with path",
			pattern: "google.com",
			url:     "https://google.com/search?q=test",
			want:    true,
		},
		{
			name:    "exact domain no match",
			pattern: "google.com",
			url:     "https://mail.google.com",
			want:    false,
		},

		// Wildcard subdomain matches
		{
			name:    "wildcard matches subdomain",
			pattern: "*.google.com",
			url:     "https://mail.google.com",
			want:    true,
		},
		{
			name:    "wildcard matches multiple subdomains",
			pattern: "*.google.com",
			url:     "https://drive.mail.google.com",
			want:    true,
		},
		{
			name:    "wildcard does not match exact domain",
			pattern: "*.google.com",
			url:     "https://google.com",
			want:    false,
		},

		// Protocol matching
		{
			name:    "https only pattern matches https",
			pattern: "https://google.com",
			url:     "https://google.com",
			want:    true,
		},
		{
			name:    "https only pattern does not match http",
			pattern: "https://google.com",
			url:     "http://google.com",
			want:    false,
		},
		{
			name:    "http only pattern matches http",
			pattern: "http://localhost",
			url:     "http://localhost",
			want:    true,
		},
		{
			name:    "http only pattern does not match https",
			pattern: "http://localhost",
			url:     "https://localhost",
			want:    false,
		},
		{
			name:    "no protocol pattern matches any protocol",
			pattern: "google.com",
			url:     "https://google.com",
			want:    true,
		},
		{
			name:    "no protocol pattern matches http",
			pattern: "google.com",
			url:     "http://google.com",
			want:    true,
		},

		// Port matching
		{
			name:    "specific port matches",
			pattern: "localhost:3000",
			url:     "http://localhost:3000",
			want:    true,
		},
		{
			name:    "specific port does not match different port",
			pattern: "localhost:3000",
			url:     "http://localhost:8080",
			want:    false,
		},
		{
			name:    "specific port does not match default port",
			pattern: "localhost:3000",
			url:     "http://localhost",
			want:    false,
		},
		{
			name:    "no port pattern matches any port",
			pattern: "localhost",
			url:     "http://localhost:3000",
			want:    true,
		},
		{
			name:    "no port pattern matches default port",
			pattern: "localhost",
			url:     "http://localhost",
			want:    true,
		},
		{
			name:    "protocol and port together",
			pattern: "http://localhost:3000",
			url:     "http://localhost:3000",
			want:    true,
		},
		{
			name:    "protocol and port mismatch",
			pattern: "https://localhost:3000",
			url:     "http://localhost:3000",
			want:    false,
		},

		// Path prefix matching
		{
			name:    "path prefix matches",
			pattern: "github.com/user",
			url:     "https://github.com/user/repo",
			want:    true,
		},
		{
			name:    "path prefix matches exact",
			pattern: "github.com/user/repo",
			url:     "https://github.com/user/repo",
			want:    true,
		},
		{
			name:    "path prefix does not match different path",
			pattern: "github.com/user",
			url:     "https://github.com/settings",
			want:    false,
		},
		{
			name:    "path prefix with trailing slash",
			pattern: "github.com/user/",
			url:     "https://github.com/user/repo",
			want:    true,
		},
		{
			name:    "no path pattern matches any path",
			pattern: "github.com",
			url:     "https://github.com/user/repo",
			want:    true,
		},
		{
			name:    "path with wildcard",
			pattern: "github.com/*",
			url:     "https://github.com/user/repo",
			want:    true,
		},
		{
			name:    "root path pattern",
			pattern: "github.com/",
			url:     "https://github.com/",
			want:    true,
		},
		{
			name:    "root path pattern matches with path",
			pattern: "github.com/",
			url:     "https://github.com/user",
			want:    true,
		},

		// Combined matching
		{
			name:    "protocol, port, and path",
			pattern: "http://localhost:3000/api",
			url:     "http://localhost:3000/api/users",
			want:    true,
		},
		{
			name:    "protocol, port, path - port mismatch",
			pattern: "http://localhost:3000/api",
			url:     "http://localhost:8080/api/users",
			want:    false,
		},
		{
			name:    "protocol, port, path - path mismatch",
			pattern: "http://localhost:3000/api",
			url:     "http://localhost:3000/web/users",
			want:    false,
		},
		{
			name:    "wildcard subdomain with path",
			pattern: "*.google.com/search",
			url:     "https://www.google.com/search?q=test",
			want:    true,
		},

		// Case insensitivity
		{
			name:    "case insensitive domain",
			pattern: "google.com",
			url:     "https://Google.COM",
			want:    true,
		},
		{
			name:    "case insensitive protocol",
			pattern: "HTTPS://google.com",
			url:     "https://google.com",
			want:    true,
		},
		{
			name:    "case sensitive path",
			pattern: "github.com/User",
			url:     "https://github.com/user",
			want:    false,
		},

		// Wildcard matches any
		{
			name:    "single wildcard matches anything",
			pattern: "*",
			url:     "https://google.com/search",
			want:    true,
		},

		// Edge cases
		{
			name:    "empty pattern",
			pattern: "",
			url:     "https://google.com",
			want:    false,
		},
		{
			name:    "pattern with www",
			pattern: "www.google.com",
			url:     "https://www.google.com",
			want:    true,
		},
		{
			name:    "url with query string",
			pattern: "google.com/search",
			url:     "https://google.com/search?q=test&lang=en",
			want:    true,
		},
		{
			name:    "url with fragment",
			pattern: "github.com/user/repo",
			url:     "https://github.com/user/repo#readme",
			want:    true,
		},

		// IP addresses
		{
			name:    "localhost ip",
			pattern: "127.0.0.1",
			url:     "http://127.0.0.1",
			want:    true,
		},
		{
			name:    "ip with port",
			pattern: "192.168.1.1:8080",
			url:     "http://192.168.1.1:8080",
			want:    true,
		},
		{
			name:    "ip with path",
			pattern: "127.0.0.1/admin",
			url:     "http://127.0.0.1/admin/users",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Match(tt.pattern, tt.url)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.url, got, tt.want)
			}
		})
	}
}

func TestMatchInvalid(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		url     string
	}{
		{
			name:    "invalid url",
			pattern: "google.com",
			url:     "not a url",
		},
		{
			name:    "invalid url scheme",
			pattern: "google.com",
			url:     "://google.com",
		},
		{
			name:    "empty url",
			pattern: "google.com",
			url:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Match(tt.pattern, tt.url)
			if got {
				t.Errorf("Match(%q, %q) = true, want false for invalid URL", tt.pattern, tt.url)
			}
		})
	}
}

func TestParsePattern(t *testing.T) {
	tests := []struct {
		name       string
		pattern    string
		wantScheme string
		wantHost   string
		wantPort   string
		wantPath   string
		wantErr    bool
	}{
		{
			name:       "simple domain",
			pattern:    "google.com",
			wantScheme: "",
			wantHost:   "google.com",
			wantPort:   "",
			wantPath:   "",
			wantErr:    false,
		},
		{
			name:       "domain with protocol",
			pattern:    "https://google.com",
			wantScheme: "https",
			wantHost:   "google.com",
			wantPort:   "",
			wantPath:   "",
			wantErr:    false,
		},
		{
			name:       "domain with port",
			pattern:    "localhost:3000",
			wantScheme: "",
			wantHost:   "localhost",
			wantPort:   "3000",
			wantPath:   "",
			wantErr:    false,
		},
		{
			name:       "domain with path",
			pattern:    "github.com/user/repo",
			wantScheme: "",
			wantHost:   "github.com",
			wantPort:   "",
			wantPath:   "/user/repo",
			wantErr:    false,
		},
		{
			name:       "full pattern",
			pattern:    "http://localhost:3000/api/users",
			wantScheme: "http",
			wantHost:   "localhost",
			wantPort:   "3000",
			wantPath:   "/api/users",
			wantErr:    false,
		},
		{
			name:       "wildcard subdomain",
			pattern:    "*.google.com",
			wantScheme: "",
			wantHost:   "*.google.com",
			wantPort:   "",
			wantPath:   "",
			wantErr:    false,
		},
		{
			name:       "wildcard subdomain with path",
			pattern:    "*.google.com/search",
			wantScheme: "",
			wantHost:   "*.google.com",
			wantPort:   "",
			wantPath:   "/search",
			wantErr:    false,
		},
		{
			name:       "single wildcard",
			pattern:    "*",
			wantScheme: "",
			wantHost:   "*",
			wantPort:   "",
			wantPath:   "",
			wantErr:    false,
		},
		{
			name:       "trailing slash",
			pattern:    "github.com/",
			wantScheme: "",
			wantHost:   "github.com",
			wantPort:   "",
			wantPath:   "/",
			wantErr:    false,
		},
		{
			name:       "path with trailing slash",
			pattern:    "github.com/user/",
			wantScheme: "",
			wantHost:   "github.com",
			wantPort:   "",
			wantPath:   "/user/",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScheme, gotHost, gotPort, gotPath, err := parsePattern(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePattern(%q) error = %v, wantErr %v", tt.pattern, err, tt.wantErr)
				return
			}
			if gotScheme != tt.wantScheme {
				t.Errorf("parsePattern(%q) scheme = %q, want %q", tt.pattern, gotScheme, tt.wantScheme)
			}
			if gotHost != tt.wantHost {
				t.Errorf("parsePattern(%q) host = %q, want %q", tt.pattern, gotHost, tt.wantHost)
			}
			if gotPort != tt.wantPort {
				t.Errorf("parsePattern(%q) port = %q, want %q", tt.pattern, gotPort, tt.wantPort)
			}
			if gotPath != tt.wantPath {
				t.Errorf("parsePattern(%q) path = %q, want %q", tt.pattern, gotPath, tt.wantPath)
			}
		})
	}
}
