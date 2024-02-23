package util

// ConvertTempToFahrenheit converts temperature from Celsius / Kelvin to Fahrenheit.
func ConvertTempToFahrenheit(temp float64, unit string) float64 {
	switch unit {
	case "metric": // Celsius to Fahrenheit
		return (temp * 9 / 5) + 32
	case "standard": // Kelvin to Fahrenheit
		return ((temp - 273.15) * 9 / 5) + 32
	default: // Assume already in Fahrenheit or an unrecognized unit, return as is
		return temp
	}
}
