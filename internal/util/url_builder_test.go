package util

import (
	"testing"
)

func TestBuildOpenWeatherMapURL(t *testing.T) {
	baseURL := "https://api.openweathermap.org/data/2.5/weather"
	apiKey := "testapikey"
	lat := "35.6895"
	lon := "139.6917"
	unitOfMeasurement := "metric"

	expectedURL := "https://api.openweathermap.org/data/2.5/weather?appid=testapikey&lat=35.6895&lon=139.6917&units=metric"

	generatedURL, err := BuildOpenWeatherMapURL(baseURL, apiKey, lat, lon, unitOfMeasurement)
	if err != nil {
		t.Fatalf("BuildOpenWeatherMapURL returned an unexpected error: %v", err)
	}

	if generatedURL != expectedURL {
		t.Errorf("Expected URL to be %v, got %v", expectedURL, generatedURL)
	}
}

func TestBuildOpenWeatherMapURLWithError(t *testing.T) {
	baseURL := "http://[::1]:namedport"
	apiKey := "testapikey"
	lat := "35.6895"
	lon := "139.6917"
	unitOfMeasurement := "metric"

	_, err := BuildOpenWeatherMapURL(baseURL, apiKey, lat, lon, unitOfMeasurement)
	if err == nil {
		t.Fatal("Expected error for invalid baseURL, but got nil")
	}
}
