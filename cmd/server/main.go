package main

import (
	"flag"
	"log"
	"spiderlite/internal/database"
	"spiderlite/internal/server"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP server address")
	dbPath := flag.String("db", "crawler.db", "Path to SQLite database")
	flag.Parse()

	// Initialize database
	db, err := database.NewDB(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create and start server
	srv := server.New(db)
	log.Printf("Starting server on %s", *addr)
	if err := srv.Start(*addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
