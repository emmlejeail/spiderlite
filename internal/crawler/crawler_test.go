package crawler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"spiderlite/internal/database"
	"spiderlite/internal/metrics"
)

func TestCrawler(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			w.Write([]byte(`
				<html><body>
					<a href="/page1">Page 1</a>
					<a href="/page2">Page 2</a>
				</body></html>
			`))
		case "/page1", "/page2":
			w.Write([]byte(`<html><body>Test page</body></html>`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	// Create test database
	db, err := database.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Use mock metrics
	m := metrics.NewNoopMetrics()

	// Create crawler
	c := New(db, m)

	// Start crawl
	startURL, _ := url.Parse(ts.URL)
	if err := c.Start(startURL); err != nil {
		t.Errorf("Crawl failed: %v", err)
	}

	// Verify results
	pages, err := db.GetPages()
	if err != nil {
		t.Errorf("Failed to get pages: %v", err)
	}

	// Should have crawled 3 pages (/, /page1, /page2)
	if len(pages) != 3 {
		t.Errorf("Expected 3 pages, got %d", len(pages))
	}

	// All pages should have status 200
	for _, page := range pages {
		if page.StatusCode != 200 {
			t.Errorf("Expected status 200 for %s, got %d", page.URL, page.StatusCode)
		}
	}
}
