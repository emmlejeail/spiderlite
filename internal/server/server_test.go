package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"spiderlite/internal/database"
	"spiderlite/internal/metrics"
)

func TestServer(t *testing.T) {
	// Create test database
	db, err := database.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Use mock metrics instead of real Datadog client
	m := metrics.NewNoopMetrics()

	// Create server
	srv := New(db, m)

	// Test cases
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "get pages",
			method:     "GET",
			path:       "/pages",
			wantStatus: http.StatusOK,
		},
		{
			name:       "get pages wrong method",
			method:     "POST",
			path:       "/pages",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "crawl without url",
			method:     "POST",
			path:       "/crawl",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			switch tt.path {
			case "/pages":
				srv.handleGetPages(w, req)
			case "/crawl":
				srv.handleCrawl(w, req)
			}

			if w.Code != tt.wantStatus {
				t.Errorf("Want status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
