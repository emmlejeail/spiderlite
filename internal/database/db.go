package database

import (
	"database/sql"
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
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
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
	query := `
	INSERT OR REPLACE INTO pages (url, status_code, crawled_at)
	VALUES (?, ?, ?)`

	_, err := db.Exec(query, page.URL, page.StatusCode, page.CrawledAt)
	return err
}

func (db *DB) GetPages() ([]PageData, error) {
	rows, err := db.Query(`
		SELECT url, status_code, crawled_at
		FROM pages
		ORDER BY crawled_at DESC
	`)
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

func (db *DB) GetPagesByStatus(status string) ([]PageData, error) {
	rows, err := db.Query(`
		SELECT url, status_code, crawled_at
		FROM pages
		WHERE status_code = ?
		ORDER BY crawled_at DESC
	`, status)
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
