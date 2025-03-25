package parser

import (
	"net/url"
	"strings"
	"testing"
)

func TestExtractLinks(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		baseURL  string
		expected []string
		wantErr  bool
	}{
		{
			name: "basic links",
			html: `<html><body>
				<a href="/page1">Page 1</a>
				<a href="https://example.com/page2">Page 2</a>
				</body></html>`,
			baseURL:  "https://example.com",
			expected: []string{"https://example.com/page1", "https://example.com/page2"},
		},
		{
			name: "invalid links",
			html: `<html><body>
				<a href=":invalid">Invalid</a>
				<a href="https://example.com/valid">Valid</a>
				</body></html>`,
			baseURL:  "https://example.com",
			expected: []string{"https://example.com/valid"},
		},
		{
			name:     "empty html",
			html:     "",
			baseURL:  "https://example.com",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseURL, err := url.Parse(tt.baseURL)
			if err != nil {
				t.Fatalf("Failed to parse base URL: %v", err)
			}

			links, err := ExtractLinks(strings.NewReader(tt.html), baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Convert []*url.URL to []string for easier comparison
			got := make([]string, len(links))
			for i, link := range links {
				got[i] = link.String()
			}

			// Compare results
			if len(got) != len(tt.expected) {
				t.Errorf("Expected %d links, got %d", len(tt.expected), len(got))
				return
			}

			for i, want := range tt.expected {
				if got[i] != want {
					t.Errorf("Link %d: want %s, got %s", i, want, got[i])
				}
			}
		})
	}
}
