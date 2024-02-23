package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang2go/demo-app/weather-service-api/internal/config"
	"github.com/golang2go/demo-app/weather-service-api/internal/model"
	"github.com/golang2go/demo-app/weather-service-api/internal/repo"
	"github.com/stretchr/testify/assert"
)

// MockWeatherAPI implementation for testing
type MockWeatherAPI struct {
	FetchFunc func(ctx context.Context, lat, lon, openWeatherMapAPIURL, unitsOfMeasurement string) (model.WeatherData, error)
}

func (m *MockWeatherAPI) FetchWeatherData(ctx context.Context, lat, lon, openWeatherMapAPIURL, unitsOfMeasurement string) (model.WeatherData, error) {
	return m.FetchFunc(ctx, lat, lon, openWeatherMapAPIURL, unitsOfMeasurement)
}

func TestWeatherHandler_GetWeatherConditionByCoordinates_FetchError(t *testing.T) {
	cfg := &config.AppConfig{
		OpenWeatherMapAPIURL: "http://example.com",
		UnitOfMeasurement:    "standard",
	}

	mockAPI := &MockWeatherAPI{
		FetchFunc: func(ctx context.Context, lat, lon, openWeatherMapAPIURL, unitsOfMeasurement string) (model.WeatherData, error) {
			return model.WeatherData{}, repo.ErrServiceUnavailable // Use the appropriate error
		},
	}

	h := NewWeatherHandler(mockAPI, cfg)

	req, _ := http.NewRequest("GET", "/weather?lat=35&lon=139", nil)
	rr := httptest.NewRecorder()
	h.GetWeatherConditionByCoordinates(rr, req)

	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
	assert.Contains(t, rr.Body.String(), "OpenWeather API service is unavailable")
}

func TestWeatherHandler_GetWeatherConditionByCoordinates_MissingParams(t *testing.T) {
	cfg := &config.AppConfig{
		OpenWeatherMapAPIURL: "http://example.com",
		UnitOfMeasurement:    "standard",
	}

	mockAPI := &MockWeatherAPI{} // No FetchFunc needed as it should not be called
	h := NewWeatherHandler(mockAPI, cfg)

	testCases := []struct {
		name       string
		query      string
		wantStatus int
		wantBody   string
	}{
		{"Missing Lat", "/weather?lon=139", http.StatusBadRequest, "Missing required query parameters: lat and/or lon\n"},
		{"Missing Lon", "/weather?lat=35", http.StatusBadRequest, "Missing required query parameters: lat and/or lon\n"},
		{"Missing Both", "/weather", http.StatusBadRequest, "Missing required query parameters: lat and/or lon\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tc.query, nil)
			rr := httptest.NewRecorder()
			h.GetWeatherConditionByCoordinates(rr, req)
			assert.Equal(t, tc.wantStatus, rr.Code)
			assert.Equal(t, tc.wantBody, rr.Body.String())
		})
	}
}

func TestWeatherHandler_GetWeatherConditionByCoordinates_Success(t *testing.T) {
	// Mock the weather API response
	mockAPI := &MockWeatherAPI{
		FetchFunc: func(ctx context.Context, lat, lon, openWeatherMapAPIURL, unitsOfMeasurement string) (model.WeatherData, error) {
			return model.WeatherData{
				Main: model.MainInfo{Temp: 280.32}, // Assuming this temperature is in Kelvin
				Weather: []model.WeatherCondition{
					{Main: "Clear"},
				},
			}, nil
		},
	}

	// Mock AppConfig for testing
	cfg := &config.AppConfig{
		OpenWeatherMapAPIURL: "http://example.com",
		UnitOfMeasurement:    "standard",
	}

	// Create an instance of WeatherHandler with mock dependencies
	h := NewWeatherHandler(mockAPI, cfg)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/weather?lat=35&lon=139", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	h.GetWeatherConditionByCoordinates(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	expected := `{"weatherCondition":"Clear","tempCategory":"Cold"}`
	assert.JSONEq(t, expected, rr.Body.String(), "handler returned unexpected body")

}

func TestCategorizeTemperature(t *testing.T) {
	tests := []struct {
		name             string
		tempFahrenheit   float64
		expectedCategory string
	}{
		{"Freezing", 32, "Freezing"},
		{"Cold", 45, "Cold"},
		{"Cool", 60, "Cool"},
		{"Mild", 70, "Mild"},
		{"Warm", 85, "Warm"},
		{"Hot", 96, "Hot"},
		{"Below Freezing", 0, "Freezing"},
		{"Above Hot", 100, "Hot"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := CategorizeTemperature(tt.tempFahrenheit)
			if category != tt.expectedCategory {
				t.Errorf("%s: CategorizeTemperature(%f) = %s, want %s", tt.name, tt.tempFahrenheit, category, tt.expectedCategory)
			}
		})
	}
}

func TestMapWeatherDataToResponse(t *testing.T) {
	tests := []struct {
		name                 string
		data                 model.WeatherData
		unitOfMeasurement    string
		expectedCondition    string
		expectedTempCategory string
	}{
		{
			name: "Clear and Freezing",
			data: model.WeatherData{
				Main:    model.MainInfo{Temp: -10}, // Celsius, translates to 14 Fahrenheit, which is Freezing
				Weather: []model.WeatherCondition{{Main: "Clear"}},
			},
			unitOfMeasurement:    "metric",
			expectedCondition:    "Clear",
			expectedTempCategory: "Freezing",
		},
		{
			name: "Rainy and Mild",
			data: model.WeatherData{
				Main:    model.MainInfo{Temp: 75}, // Fahrenheit, directly Mild
				Weather: []model.WeatherCondition{{Main: "Rain"}},
			},
			unitOfMeasurement:    "imperial",
			expectedCondition:    "Rain",
			expectedTempCategory: "Mild",
		},
		// other usecases with different unit of measurements / conditions / temp category
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := MapWeatherDataToResponse(tt.data, tt.unitOfMeasurement)

			if response.WeatherCondition != tt.expectedCondition {
				t.Errorf("%s: expected condition %s, got %s", tt.name, tt.expectedCondition, response.WeatherCondition)
			}

			if response.TempCategory != tt.expectedTempCategory {
				t.Errorf("%s: expected temperature category %s, got %s", tt.name, tt.expectedTempCategory, response.TempCategory)
			}
		})
	}
}

func TestHandleWeatherDataError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Invalid API Key",
			err:            repo.ErrInvalidAPIKey,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid API key.\n",
		},
		{
			name:           "Bad Request",
			err:            repo.ErrBadRequest,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad request to OpenWeather API.\n",
		},
		{
			name:           "Service Unavailable",
			err:            repo.ErrServiceUnavailable,
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   "OpenWeather API service is unavailable.\n",
		},
		{
			name:           "Unexpected Status Code",
			err:            repo.ErrUnexpectedStatusCode,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "An error occurred while processing your request.\n",
		},
		{
			name:           "Decoding Response Error",
			err:            repo.ErrDecodingResponse,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "An error occurred while processing your request.\n",
		},
		{
			name:           "Unknown Error",
			err:            errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "An unexpected error occurred.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			recorder := httptest.NewRecorder()

			// Act
			handleWeatherDataError(tt.err, recorder)

			// Assert
			result := recorder.Result()
			defer result.Body.Close()

			// Check status code
			assert.Equal(t, tt.expectedStatus, result.StatusCode, tt.name+" - StatusCode")

			// Read and check the body
			bodyBytes, err := io.ReadAll(result.Body)
			assert.NoError(t, err, tt.name+" - Reading body")
			body := string(bodyBytes)
			assert.Equal(t, tt.expectedBody, body, tt.name+" - Body")
		})
	}
}
