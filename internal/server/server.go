package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"spiderlite/internal/crawler"
	"spiderlite/internal/database"
	"spiderlite/internal/metrics"
	"strconv"
)

type Server struct {
	db      *database.DB
	metrics *metrics.Metrics
	crawler *crawler.Crawler
}

func New(db *database.DB, metrics *metrics.Metrics) *Server {
	return &Server{
		db:      db,
		metrics: metrics,
		crawler: crawler.New(db, metrics),
	}
}

func (s *Server) handleCrawl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get URL from query parameters
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	log.Printf("Received crawl request for URL: %s", targetURL)

	// Validate URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		log.Printf("Invalid URL: %v", err)
		http.Error(w, "Invalid URL: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Launch crawl in a goroutine
	go func() {
		log.Printf("Starting crawl for: %s", parsedURL.String())
		if err := s.crawler.Start(parsedURL); err != nil {
			log.Printf("Crawl error: %v", err)
			s.metrics.IncrementCrawlErrors()
		}
	}()

	response := map[string]string{
		"status":  "started",
		"message": "Crawl started for " + targetURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleDebug(w http.ResponseWriter, r *http.Request) {
	var debug struct {
		DatabasePath string           `json:"database_path"`
		TableCount   int              `json:"table_count"`
		Tables       map[string]int64 `json:"tables"`
		Error        string           `json:"error,omitempty"`
	}

	debug.Tables = make(map[string]int64)

	// Get database path
	var seq int
	var name, file string
	err := s.db.QueryRow("PRAGMA database_list").Scan(&seq, &name, &file)
	if err != nil {
		debug.Error = "Failed to get DB path: " + err.Error()
	}
	debug.DatabasePath = file

	// Get table info
	row := s.db.QueryRow("SELECT COUNT(*) FROM pages")
	if err := row.Scan(&debug.TableCount); err != nil {
		debug.Error = fmt.Sprintf("%s; Failed to get count: %v", debug.Error, err)
	}

	// Get table list
	rows, err := s.db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err == nil {
				var count int64
				if err := s.db.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count); err == nil {
					debug.Tables[tableName] = count
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(debug)
}

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/pages", metricsMiddleware(s.metrics, "/pages")(s.handleGetPages))
	mux.HandleFunc("/pages/status", metricsMiddleware(s.metrics, "/pages/status")(s.handleGetPagesByStatus))
	mux.HandleFunc("/crawl", metricsMiddleware(s.metrics, "/crawl")(s.handleCrawl))
	mux.HandleFunc("/debug", s.handleDebug)

	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleGetPages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Fetching pages from database...")
	pages, err := s.db.GetPages()
	if err != nil {
		log.Printf("Error fetching pages: %v", err)
		http.Error(w, "Failed to fetch pages: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Found %d pages in database", len(pages))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(pages),
		"pages": pages,
	})
}

func (s *Server) handleGetPagesByStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	statusStr := r.URL.Query().Get("code")
	if statusStr == "" {
		http.Error(w, "Status code required", http.StatusBadRequest)
		return
	}

	statusCode, err := strconv.Atoi(statusStr)
	if err != nil {
		http.Error(w, "Invalid status code", http.StatusBadRequest)
		return
	}

	pages, err := s.db.GetPagesByStatus(statusCode)
	if err != nil {
		http.Error(w, "Failed to fetch pages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": statusCode,
		"count":  len(pages),
		"pages":  pages,
	})
}
