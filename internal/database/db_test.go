package database

import (
	"testing"
	"time"
)

func TestDatabase(t *testing.T) {
	// Use in-memory SQLite for testing
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test data
	testPages := []PageData{
		{
			URL:        "https://example.com",
			StatusCode: 200,
			CrawledAt:  time.Now(),
		},
		{
			URL:        "https://example.com/page1",
			StatusCode: 404,
			CrawledAt:  time.Now(),
		},
	}

	// Test storing pages
	t.Run("store pages", func(t *testing.T) {
		for _, page := range testPages {
			if err := db.StorePage(page); err != nil {
				t.Errorf("StorePage() error = %v", err)
			}
		}
	})

	// Test retrieving all pages
	t.Run("get all pages", func(t *testing.T) {
		pages, err := db.GetPages()
		if err != nil {
			t.Errorf("GetPages() error = %v", err)
			return
		}

		if len(pages) != len(testPages) {
			t.Errorf("Expected %d pages, got %d", len(testPages), len(pages))
		}
	})

	// Test retrieving pages by status
	t.Run("get pages by status", func(t *testing.T) {
		pages, err := db.GetPagesByStatus(200)
		if err != nil {
			t.Errorf("GetPagesByStatus() error = %v", err)
			return
		}

		if len(pages) != 1 {
			t.Errorf("Expected 1 page with status 200, got %d", len(pages))
		}
	})
}
