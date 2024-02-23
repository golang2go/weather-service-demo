package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// Create a test handler that responds with a specific status code
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound) // Example status code to test logging of non-200 responses
	})

	// Wrap the testHandler with the LoggingMiddleware
	middlewareHandler := LoggingMiddleware(testHandler)

	// Create a ResponseRecorder (to capture the response) and a dummy request
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/testpath", nil)

	// Serve the request through the middleware
	middlewareHandler.ServeHTTP(rr, req)

	// Check the status code to verify the middleware correctly allowed the testHandler to set it
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	// Note: This test does not verify the log output, as capturing log output from the log package requires
	// redirecting the output to a buffer and parsing it, which is beyond the scope of this basic test.
}
