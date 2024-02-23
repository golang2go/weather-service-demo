package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimitMiddleware(t *testing.T) {
	// Create a simple test handler that returns 200 OK
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap testHandler with the RateLimitMiddleware set to allow 1 request per second
	middlewareHandler := RateLimitMiddleware(1)(testHandler)

	// Create a ResponseRecorder (to record responses) and a dummy request
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	// First request should pass through
	middlewareHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("First request was unexpectedly blocked by rate limiting")
	}

	// Immediate subsequent request should be rate limited
	rr2 := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(rr2, req)

	if status := rr2.Code; status != http.StatusTooManyRequests {
		t.Errorf("Second immediate request was not rate limited as expected")
	}

	// Check for the presence of the Retry-After header
	if retryAfter := rr2.Header().Get("Retry-After"); retryAfter == "" {
		t.Errorf("Rate limited response did not include a Retry-After header")
	}

	// Wait for the rate limit duration plus a small buffer to ensure the rate limiter has reset
	time.Sleep(time.Second + 10*time.Millisecond)

	// Next request after waiting should pass through
	rr3 := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(rr3, req)

	if status := rr3.Code; status != http.StatusOK {
		t.Errorf("Request after rate limit duration was unexpectedly blocked")
	}
}
