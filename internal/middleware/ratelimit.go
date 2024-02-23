package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// RateLimitMiddleware limits the number of requests handled by the server based on a configurable rate limit.
func RateLimitMiddleware(rateLimitPerSecond int) func(http.Handler) http.Handler {
	// Define a mutex to protect lastRequestTime
	var mutex sync.Mutex
	// Initialize lastRequestTime to the current time
	lastRequestTime := time.Now()

	// Calculate the limit duration based on the rate limit per second
	limitDuration := time.Second / time.Duration(rateLimitPerSecond)

	// Return the middleware handler function
	return func(next http.Handler) http.Handler {
		// Middleware handler function
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Acquire the mutex lock before accessing lastRequestTime
			mutex.Lock()
			defer mutex.Unlock() // Ensure the mutex is always released

			// Calculate the time since the last request
			timeSinceLastRequest := time.Since(lastRequestTime)

			// If the time since the last request is less than the limit duration,
			// return a rate limit exceeded error response
			if timeSinceLastRequest < limitDuration {
				// Calculate the retry after duration
				retryAfter := limitDuration - timeSinceLastRequest
				// Set the Retry-After header in the response
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(retryAfter.Seconds())+1))
				// Return a rate limit exceeded error response
				http.Error(w, fmt.Sprintf("Rate limit exceeded. Please wait %d seconds before retrying.", int(retryAfter.Seconds())+1), http.StatusTooManyRequests)
				return
			}

			// Update the lastRequestTime to the current time
			lastRequestTime = time.Now()

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
