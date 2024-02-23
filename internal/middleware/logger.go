package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter creates a new response writer to capture the status code.
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	// Default the status code to 200 for cases where WriteHeader is not explicitly called.
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer
		wrappedWriter := newResponseWriter(w)

		// Process the request
		next.ServeHTTP(wrappedWriter, r)

		// Log the request details including the status code
		log.Printf("Method: %s, URI: %s, Status: %d, Duration: %v",
			r.Method, r.URL.Path, wrappedWriter.statusCode, time.Since(start))
	})
}
