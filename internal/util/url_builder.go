package util

import (
	"fmt"
	"net/url"
)

func BuildOpenWeatherMapURL(baseURL, apiKey, lat, lon, unitOfMeasurement string) (string, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("error parsing base URL: %v", err)
	}

	query := parsedURL.Query()
	query.Set("lat", lat)
	query.Set("lon", lon)
	query.Set("appid", apiKey)
	query.Set("units", unitOfMeasurement)

	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}
