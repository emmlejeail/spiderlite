package server

import (
	"encoding/json"
	"net/http"
	"spiderlite/internal/database"
	"spiderlite/internal/metrics"
)

type Server struct {
	db      *database.DB
	metrics *metrics.Metrics
}

func New(db *database.DB, metrics *metrics.Metrics) *Server {
	return &Server{
		db:      db,
		metrics: metrics,
	}
}

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()

	// Register routes with metrics middleware
	mux.HandleFunc("/pages", metricsMiddleware(s.metrics, "/pages")(s.handleGetPages))
	mux.HandleFunc("/pages/status", metricsMiddleware(s.metrics, "/pages/status")(s.handleGetPagesByStatus))

	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleGetPages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pages, err := s.db.GetPages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pages)
}

func (s *Server) handleGetPagesByStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := r.URL.Query().Get("code")
	if status == "" {
		http.Error(w, "Status code required", http.StatusBadRequest)
		return
	}

	pages, err := s.db.GetPagesByStatus(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pages)
}
