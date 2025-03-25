package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

type PageData struct {
	URL        string
	StatusCode int
	CrawledAt  time.Time
}

func NewDB(dbPath string) (*DB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable foreign keys and WAL mode
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, err
	}

	if err := initSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db}, nil
}

func initSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS pages (
		url TEXT PRIMARY KEY,
		status_code INTEGER,
		crawled_at DATETIME
	);`

	_, err := db.Exec(query)
	return err
}

func (db *DB) StorePage(page PageData) error {
	log.Printf("Attempting to store page: %s", page.URL)

	query := `
	INSERT OR REPLACE INTO pages (url, status_code, crawled_at)
	VALUES (?, ?, ?)`

	result, err := db.Exec(query, page.URL, page.StatusCode, page.CrawledAt)
	if err != nil {
		log.Printf("Error storing page: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Successfully stored page %s. Rows affected: %d", page.URL, rowsAffected)
	return nil
}

func (db *DB) GetPages() ([]PageData, error) {
	log.Printf("Executing GetPages query...")

	query := `
		SELECT url, status_code, crawled_at
		FROM pages
		ORDER BY crawled_at DESC
		LIMIT 100`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying pages: %v", err)
		return nil, err
	}
	defer rows.Close()

	var pages []PageData
	for rows.Next() {
		var page PageData
		err := rows.Scan(&page.URL, &page.StatusCode, &page.CrawledAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		pages = append(pages, page)
	}

	log.Printf("Retrieved %d pages from database", len(pages))
	return pages, nil
}

func (db *DB) GetPagesByStatus(statusCode int) ([]PageData, error) {
	query := `
		SELECT url, status_code, crawled_at
		FROM pages
		WHERE status_code = ?
		ORDER BY crawled_at DESC
		LIMIT 100`

	rows, err := db.Query(query, statusCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []PageData
	for rows.Next() {
		var page PageData
		err := rows.Scan(&page.URL, &page.StatusCode, &page.CrawledAt)
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	return pages, nil
}
