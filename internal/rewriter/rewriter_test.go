package rewriter

import (
	"testing"
)

func TestRewrite(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		template string
		want     string
		wantErr  bool
	}{
		{
			name:     "empty template returns original URL",
			url:      "https://example.com/path",
			template: "",
			want:     "https://example.com/path",
			wantErr:  false,
		},
		{
			name:     "simple host replacement",
			url:      "https://twitter.com/user/status/123",
			template: "xcancel.com{path}",
			want:     "https://xcancel.com/user/status/123",
			wantErr:  false,
		},
		{
			name:     "full URL with all components",
			url:      "https://example.com:8080/path/to/page?key=value&key2=value2#section",
			template: "newhost.com{path}?{query}#{fragment}",
			want:     "https://newhost.com/path/to/page?key=value&key2=value2#section",
			wantErr:  false,
		},
		{
			name:     "preserve scheme",
			url:      "http://example.com/path",
			template: "newhost.com{path}",
			want:     "http://newhost.com/path",
			wantErr:  false,
		},
		{
			name:     "use original host in template",
			url:      "https://sub.example.com/path",
			template: "archive.{host}{path}",
			want:     "https://archive.sub.example.com/path",
			wantErr:  false,
		},
		{
			name:     "youtube to invidious",
			url:      "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			template: "invidious.io{path}?{query}",
			want:     "https://invidious.io/watch?v=dQw4w9WgXcQ",
			wantErr:  false,
		},
		{
			name:     "reddit to teddit",
			url:      "https://old.reddit.com/r/programming/comments/abc123/title",
			template: "teddit.net{path}",
			want:     "https://teddit.net/r/programming/comments/abc123/title",
			wantErr:  false,
		},
		{
			name:     "preserve query string only",
			url:      "https://example.com/search?q=test&sort=date",
			template: "searx.example.com/search?{query}",
			want:     "https://searx.example.com/search?q=test&sort=date",
			wantErr:  false,
		},
		{
			name:     "path with trailing slash",
			url:      "https://example.com/path/",
			template: "newhost.com{path}",
			want:     "https://newhost.com/path/",
			wantErr:  false,
		},
		{
			name:     "no path in original URL",
			url:      "https://example.com",
			template: "newhost.com{path}",
			want:     "https://newhost.com",
			wantErr:  false,
		},
		{
			name:     "port in original URL",
			url:      "https://example.com:3000/path",
			template: "newhost.com:{port}{path}",
			want:     "https://newhost.com:3000/path",
			wantErr:  false,
		},
		{
			name:     "fragment without query",
			url:      "https://example.com/path#section",
			template: "newhost.com{path}#{fragment}",
			want:     "https://newhost.com/path#section",
			wantErr:  false,
		},
		{
			name:     "scheme in template",
			url:      "https://example.com/path",
			template: "{scheme}://newhost.com{path}",
			want:     "https://newhost.com/path",
			wantErr:  false,
		},
		{
			name:     "complex rewrite with all variables",
			url:      "https://old.example.com:8080/api/v1/users?limit=10&offset=20#results",
			template: "{scheme}://new.{host}:{port}/v2{path}?{query}#{fragment}",
			want:     "https://new.old.example.com:8080/v2/api/v1/users?limit=10&offset=20#results",
			wantErr:  false,
		},
		{
			name:     "URL without scheme gets https by default",
			url:      "example.com/path",
			template: "newhost.com{path}",
			want:     "https://newhost.com/path",
			wantErr:  false,
		},
		{
			name:     "empty query and fragment get normalized",
			url:      "https://example.com/path",
			template: "newhost.com{path}?{query}#{fragment}",
			want:     "https://newhost.com/path?",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Rewrite(tt.url, tt.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rewrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Rewrite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRewriteInvalidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		template string
	}{
		{
			name:     "invalid URL with spaces",
			url:      "https://example.com/path with spaces",
			template: "newhost.com{path}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Rewrite(tt.url, tt.template)
			// Should return original URL when there's an error
			if err != nil && got != tt.url {
				t.Errorf("Rewrite() with invalid URL should return original URL, got %v", got)
			}
		})
	}
}
