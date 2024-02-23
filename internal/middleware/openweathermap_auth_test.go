package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenWeatherMapAuthMiddleware(t *testing.T) {
	// Create a test handler that will be wrapped by the middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the testHandler with the OpenWeatherMapAuthMiddleware
	middlewareHandler := OpenWeatherMapAuthMiddleware(testHandler)

	t.Run("With API Key", func(t *testing.T) {
		// Create a request with the API key header
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set(APIKeyHeader, "test-api-key")

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Serve the request through the middleware
		middlewareHandler.ServeHTTP(rr, req)

		// Check if the status code is 200 OK
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Without API Key", func(t *testing.T) {
		// Create a request without the API key header
		req, _ := http.NewRequest("GET", "/", nil)

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Serve the request through the middleware
		middlewareHandler.ServeHTTP(rr, req)

		// Check if the status code is 400 Bad Request
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}
