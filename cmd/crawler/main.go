package main

import (
	"log"
	"net/url"
	"os"

	"spiderlite/internal/crawler"
	"spiderlite/internal/database"
	"spiderlite/internal/metrics"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <url>", os.Args[0])
	}
	startURL := os.Args[1]

	parsedURL, err := url.Parse(startURL)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}

	// Initialize metrics
	metrics, err := metrics.New()
	if err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}
	defer metrics.Close()

	// Initialize database
	db, err := database.NewDB("crawler.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create crawler instance with metrics
	c := crawler.New(db, metrics)

	// Start crawling
	if err := c.Start(parsedURL); err != nil {
		log.Fatalf("Crawling failed: %v", err)
	}
}
