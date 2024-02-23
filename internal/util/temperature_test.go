package util

import (
	"testing"
)

func TestConvertTempToFahrenheit(t *testing.T) {
	tests := []struct {
		name     string
		temp     float64
		unit     string
		expected float64
	}{
		{"Celsius to Fahrenheit", 0, "metric", 32},
		{"Kelvin to Fahrenheit", 273.15, "standard", 32},
		{"Already Fahrenheit", 32, "imperial", 32},
		{"Unrecognized Unit", 100, "unknown", 100},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertTempToFahrenheit(tc.temp, tc.unit)
			if result != tc.expected {
				t.Errorf("ConvertTempToFahrenheit(%f, %s) = %f; want %f", tc.temp, tc.unit, result, tc.expected)
			}
		})
	}
}
