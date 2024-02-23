package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang2go/demo-app/weather-service-api/internal/config"
	"github.com/golang2go/demo-app/weather-service-api/internal/model"
	"github.com/golang2go/demo-app/weather-service-api/internal/repo"
	"github.com/golang2go/demo-app/weather-service-api/internal/util"
)

// WeatherHandler handles weather-related HTTP requests by fetching weather data and using the application's configuration.
type WeatherHandler struct {
	Repo   repo.WeatherAPI   // Interface for fetching weather data
	Config *config.AppConfig // Application configuration settings
}

// NewWeatherHandler creates a WeatherHandler with given weather data repository and configuration for easier testing.
func NewWeatherHandler(repo repo.WeatherAPI, cfg *config.AppConfig) *WeatherHandler {
	return &WeatherHandler{Repo: repo, Config: cfg}
}

// GetWeatherConditionByCoordinates handles HTTP requests for weather conditions by coordinates.
func (h *WeatherHandler) GetWeatherConditionByCoordinates(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()
	lat := query.Get("lat")
	lon := query.Get("lon")

	// Validate query parameters
	if lat == "" || lon == "" {
		http.Error(w, "Missing required query parameters: lat and/or lon", http.StatusBadRequest)
		return
	}

	// Call the OpenWeather API using the fetcher
	weatherData, err := h.Repo.FetchWeatherData(r.Context(), lat, lon, h.Config.OpenWeatherMapAPIURL, h.Config.UnitOfMeasurement)
	if err != nil {
		handleWeatherDataError(err, w)
		return
	}

	// Map the API response to the response model
	response := MapWeatherDataToResponse(weatherData, h.Config.UnitOfMeasurement)

	// Respond to the client with the weather condition and temperature category
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// mapWeatherDataToResponse maps data from the OpenWeather API to the custom response format.
func MapWeatherDataToResponse(data model.WeatherData, unitOfMeasurement string) model.WeatherResponse {
	condition := "Unknown" // Default condition if none is found
	if len(data.Weather) > 0 {
		condition = data.Weather[0].Main
	}

	tempFahrenheit := util.ConvertTempToFahrenheit(data.Main.Temp, unitOfMeasurement)
	tempCategory := CategorizeTemperature(tempFahrenheit)

	return model.WeatherResponse{
		WeatherCondition: condition,
		TempCategory:     tempCategory,
	}
}

// CategorizeTemperature categorizes the temperature into human-readable form.
func CategorizeTemperature(tempFahrenheit float64) string {
	switch {
	case tempFahrenheit <= 32:
		return "Freezing"
	case tempFahrenheit > 32 && tempFahrenheit <= 50:
		return "Cold"
	case tempFahrenheit > 50 && tempFahrenheit <= 68:
		return "Cool"
	case tempFahrenheit > 68 && tempFahrenheit <= 77:
		return "Mild"
	case tempFahrenheit > 77 && tempFahrenheit <= 95:
		return "Warm"
	default:
		return "Hot"
	}
}

// handleWeatherDataError handles errors returned by FetchWeatherData.
func handleWeatherDataError(err error, w http.ResponseWriter) {
	var status int
	var message string

	switch {
	case errors.Is(err, repo.ErrInvalidAPIKey):
		status = http.StatusUnauthorized
		message = "Invalid API key."
	case errors.Is(err, repo.ErrBadRequest):
		status = http.StatusBadRequest
		message = "Bad request to OpenWeather API."
	case errors.Is(err, repo.ErrServiceUnavailable):
		status = http.StatusServiceUnavailable
		message = "OpenWeather API service is unavailable."
	case errors.Is(err, repo.ErrUnexpectedStatusCode), errors.Is(err, repo.ErrDecodingResponse):
		status = http.StatusInternalServerError
		message = "An error occurred while processing your request."
	default:
		status = http.StatusInternalServerError
		message = "An unexpected error occurred."
	}

	http.Error(w, message, status)
}
