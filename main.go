package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
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

// GET /v1/healthcheck
func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "status: available\n")
	app.logger.Info("healthcheck handler called")
}

// GET /v1/books
func (app *application) listBooks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "list of books (coming soon)\n")
	app.logger.Info("listBooks handler called")
}

// GET /v1/books/{id}
func (app *application) getBook(w http.ResponseWriter, r *http.Request) {

	// Extract ID from route (Go 1.22 feature)
	id := r.PathValue("id")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "get book with id: %s\n", id)

	app.logger.Info("getBook handler called", "id", id)
}

// POST /v1/books
func (app *application) createBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated) // 201
	fmt.Fprintf(w, "book created (coming soon)\n")
	app.logger.Info("createBook handler called")
}

// DELETE /v1/books/{id}
func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	w.WriteHeader(http.StatusNoContent) // 204

	// IMPORTANT: 204 responses must NOT have a body.
	app.logger.Info("deleteBook handler called", "id", id)
}

func main() {
	// Create structured logger writing to stdout
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Inject logger into application struct
	app := &application{
		logger: logger,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheck)
	mux.HandleFunc("GET /v1/books", app.listBooks)
	mux.HandleFunc("GET /v1/books/{id}", app.getBook)
	mux.HandleFunc("POST /v1/books", app.createBook)
	mux.HandleFunc("DELETE /v1/books/{id}", app.deleteBook)

	//////////////////////////////////////////////////////
	// Log Startup Message
	//////////////////////////////////////////////////////

	logger.Info("starting server", "addr", ":4000")

	//////////////////////////////////////////////////////
	// Wrap Entire Router with Middleware
	//////////////////////////////////////////////////////

	err := http.ListenAndServe(":4000", loggingMiddleware(mux))
	log.Fatal(err)
}
