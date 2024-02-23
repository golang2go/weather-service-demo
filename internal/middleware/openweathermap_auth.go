package middleware

import (
	"context"
	"net/http"
)

// APIKeyContextKey is a custom type to define the key for storing API key in request context
type APIKeyContextKey string

const (
	// APIKeyHeader is the header key for the API key
	APIKeyHeader = "X-API-Key"
)

// OpenWeatherMapAuthMiddleware checks for the presence of an OpenWeatherMap API key in the request header.
func OpenWeatherMapAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract API key from header
		apiKey := r.Header.Get(APIKeyHeader)
		if apiKey == "" {
			http.Error(w, "Missing 'X-API-Key' header. Include your OpenWeatherMap API key in the 'X-API-Key' header. See documentation for more details.", http.StatusBadRequest)

			return
		}

		// Add API key to request context
		ctx := context.WithValue(r.Context(), APIKeyContextKey("apiKey"), apiKey)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
