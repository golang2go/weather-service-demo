package repo

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang2go/demo-app/weather-service-api/internal/middleware"
	"github.com/golang2go/demo-app/weather-service-api/internal/model"
)

// setupMockServer helps in creating a mock server for testing
func setupMockServer(response string, statusCode int) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	})
	return httptest.NewServer(handler)
}

func TestFetchWeatherData_Success(t *testing.T) {
	mockResponse := `{"main":{"temp":280.32},"weather":[{"main":"Clear"}]}`
	mockServer := setupMockServer(mockResponse, http.StatusOK)
	defer mockServer.Close()

	api := NewWeatherAPI() // Correct instantiation

	ctx := context.WithValue(context.Background(), middleware.APIKeyContextKey("apiKey"), "valid-api-key")

	data, err := api.FetchWeatherData(ctx, "35", "139", mockServer.URL, "metric")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Unmarshal the expected response to compare
	var expectedData model.WeatherData
	json.Unmarshal([]byte(mockResponse), &expectedData)

	if data.Main.Temp != expectedData.Main.Temp {
		t.Errorf("Expected temperature of %v, got %v", expectedData.Main.Temp, data.Main.Temp)
	}

	if len(data.Weather) == 0 || data.Weather[0].Main != expectedData.Weather[0].Main {
		t.Errorf("Expected weather condition to be %v, got %v", expectedData.Weather[0].Main, data.Weather)
	}
}

func TestFetchWeatherData_Timeout(t *testing.T) {
	// Simulate a slow response or network issue that would cause a timeout
	mockServer := setupMockServer("", http.StatusOK) // The response body or status code is irrelevant here
	mockServer.Config.ConnState = func(conn net.Conn, state http.ConnState) {
		if state == http.StateNew {
			// Simulate a delay longer than the requestTimeout
			time.Sleep(requestTimeout + 1*time.Second)
		}
	}
	defer mockServer.Close()

	api := NewWeatherAPI()

	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.APIKeyContextKey("apiKey"), "valid-api-key")

	_, err := api.FetchWeatherData(ctx, "35", "139", mockServer.URL, "metric")
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("Expected ErrTimeout, got %v", err)
	}
}

func TestFetchWeatherData_InvalidAPIKey(t *testing.T) {
	mockResponse := `{"cod":401, "message":"Invalid API key"}`
	mockServer := setupMockServer(mockResponse, http.StatusUnauthorized)
	defer mockServer.Close()

	api := NewWeatherAPI() // Correct instantiation

	ctx := context.WithValue(context.Background(), middleware.APIKeyContextKey("apiKey"), "invalid-api-key")

	_, err := api.FetchWeatherData(ctx, "35", "139", mockServer.URL, "metric")
	if !errors.Is(err, ErrInvalidAPIKey) {
		t.Errorf("Expected ErrInvalidAPIKey, got %v", err)
	}
}

func TestFetchWeatherData_DecodingError(t *testing.T) {
	mockResponse := `{"main":{"temp":280.32},"weather":[{"main":"Clear"}` // Malformed JSON
	mockServer := setupMockServer(mockResponse, http.StatusOK)
	defer mockServer.Close()

	api := NewWeatherAPI()

	ctx := context.WithValue(context.Background(), middleware.APIKeyContextKey("apiKey"), "valid-api-key")

	_, err := api.FetchWeatherData(ctx, "35", "139", mockServer.URL, "metric")
	if !errors.Is(err, ErrDecodingResponse) {
		t.Errorf("Expected ErrDecodingResponse due to malformed JSON, got %v", err)
	}
}

func TestFetchWeatherData_BadRequest(t *testing.T) {
	mockResponse := `{"cod":400, "message":"Bad request"}`
	mockServer := setupMockServer(mockResponse, http.StatusBadRequest)
	defer mockServer.Close()

	api := NewWeatherAPI()

	ctx := context.WithValue(context.Background(), middleware.APIKeyContextKey("apiKey"), "valid-api-key")

	// We're using a valid URL, but the mock server will return a 400 status, simulating a bad request.
	_, err := api.FetchWeatherData(ctx, "invalid", "invalid", mockServer.URL, "metric")
	if !errors.Is(err, ErrBadRequest) {
		t.Errorf("Expected ErrBadRequest for new request error, got %v", err)
	}
}

func TestFetchWeatherData_HttpDoError(t *testing.T) {
	api := NewWeatherAPI()

	mockServer := setupMockServer("", http.StatusOK)
	mockServer.Close() // Close the server to simulate a network error

	ctx := context.WithValue(context.Background(), middleware.APIKeyContextKey("apiKey"), "valid-api-key")
	_, err := api.FetchWeatherData(ctx, "35", "139", mockServer.URL, "metric")
	if !errors.Is(err, ErrServiceUnavailable) {
		t.Errorf("Expected ErrServiceUnavailable for HTTP do error, got %v", err)
	}
}

func TestFetchWeatherData_MissingAPIKey(t *testing.T) {
	api := NewWeatherAPI()

	ctx := context.Background() // Context without API key
	_, err := api.FetchWeatherData(ctx, "35", "139", "http://example.com", "metric")
	if !errors.Is(err, ErrBadRequest) {
		t.Errorf("Expected ErrBadRequest for missing API key, got %v", err)
	}
}

func TestFetchWeatherData_ServiceUnavailable(t *testing.T) {
	mockResponse := `{"cod":503, "message":"Service unavailable"}`
	mockServer := setupMockServer(mockResponse, http.StatusServiceUnavailable)
	defer mockServer.Close()

	api := NewWeatherAPI()

	ctx := context.WithValue(context.Background(), middleware.APIKeyContextKey("apiKey"), "valid-api-key")

	_, err := api.FetchWeatherData(ctx, "35", "139", mockServer.URL, "metric")
	if !errors.Is(err, ErrServiceUnavailable) {
		t.Errorf("Expected ErrServiceUnavailable, got %v", err)
	}
}

func TestFetchWeatherData_InvalidJSONStructure(t *testing.T) {
	mockResponse := `{"main":{"temp":280.32},"weather":[{"main":"Clear"}}` // Missing closing ]
	mockServer := setupMockServer(mockResponse, http.StatusOK)
	defer mockServer.Close()

	api := NewWeatherAPI()

	ctx := context.WithValue(context.Background(), middleware.APIKeyContextKey("apiKey"), "valid-api-key")

	_, err := api.FetchWeatherData(ctx, "35", "139", mockServer.URL, "metric")
	if !errors.Is(err, ErrDecodingResponse) {
		t.Errorf("Expected ErrDecodingResponse due to invalid JSON structure, got %v", err)
	}
}

func TestFetchWeatherData_RateLimit(t *testing.T) {
	mockResponse := `{"cod":429, "message":"Too many requests"}`
	mockServer := setupMockServer(mockResponse, http.StatusTooManyRequests)
	defer mockServer.Close()

	api := NewWeatherAPI()

	ctx := context.WithValue(context.Background(), middleware.APIKeyContextKey("apiKey"), "valid-api-key")

	_, err := api.FetchWeatherData(ctx, "35", "139", mockServer.URL, "metric")
	if !errors.Is(err, ErrUnexpectedStatusCode) {
		t.Errorf("Expected ErrUnexpectedStatusCode for rate limit, got %v", err)
	}
}

func TestFetchWeatherData_UnexpectedStatusCode(t *testing.T) {
	mockResponse := `{"cod":418, "message":"I'm a teapot"}`
	mockServer := setupMockServer(mockResponse, http.StatusTeapot) // Using 418 I'm a teapot for testing
	defer mockServer.Close()

	api := NewWeatherAPI()

	ctx := context.WithValue(context.Background(), middleware.APIKeyContextKey("apiKey"), "valid-api-key")

	_, err := api.FetchWeatherData(ctx, "35", "139", mockServer.URL, "metric")
	if !errors.Is(err, ErrUnexpectedStatusCode) {
		t.Errorf("Expected ErrUnexpectedStatusCode, got %v", err)
	}
}
