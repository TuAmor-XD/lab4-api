package main

import (
	"log"
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter     // embed the real writer
	statusCode          int // captured status code
}

// WriteHeader intercepts the status code before forwarding it.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap original writer and preset status to 200
		rw := &responseWriter{w, http.StatusOK}

		// Call next handler
		next.ServeHTTP(rw, r)

		// Log method, path, status, and duration
		log.Printf("%s %s %d %v",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(start),
		)
	})
}

type application struct {
	logger *slog.Logger
}


