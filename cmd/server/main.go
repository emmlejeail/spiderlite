package main

import (
	"flag"
	"log"
	"os"
	"spiderlite/internal/database"
	"spiderlite/internal/metrics"
	"spiderlite/internal/server"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP server address")
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/data/crawler.db" // default path
	}

	log.Printf("Using database path: %s", dbPath)

	// Initialize database
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize metrics
	metrics, err := metrics.New()
	if err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}
	defer metrics.Close()

	// Create and start server
	srv := server.New(db, metrics)
	log.Printf("Starting server on %s", *addr)
	if err := srv.Start(*addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
