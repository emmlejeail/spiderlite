package main

import (
	"flag"
	"log"
	"spiderlite/internal/database"
	"spiderlite/internal/metrics"
	"spiderlite/internal/server"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP server address")
	dbPath := flag.String("db", "crawler.db", "Path to SQLite database")
	flag.Parse()

	// Initialize metrics
	metrics, err := metrics.New()
	if err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}
	defer metrics.Close()

	// Initialize database
	db, err := database.NewDB(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create and start server
	srv := server.New(db, metrics)
	log.Printf("Starting server on %s", *addr)
	if err := srv.Start(*addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
