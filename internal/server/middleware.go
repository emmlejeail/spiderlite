package server

import (
	"net/http"
	"time"

	"spiderlite/internal/metrics"
)

func metricsMiddleware(metrics *metrics.Metrics, endpoint string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap ResponseWriter pour capturer le status code
			rw := &responseWriter{w, http.StatusOK}

			next.ServeHTTP(rw, r)

			metrics.IncrementAPIRequests(endpoint, r.Method, rw.statusCode)
			metrics.TimeAPIRequest(endpoint, time.Since(start))
		}
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
