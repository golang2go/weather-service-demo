package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang2go/demo-app/weather-service-api/internal/middleware"
	"github.com/golang2go/demo-app/weather-service-api/internal/model"
	"github.com/golang2go/demo-app/weather-service-api/internal/util"
)

var (
	ErrInvalidAPIKey        = errors.New("invalid OpenWeatherMap API key")
	ErrBadRequest           = errors.New("bad request to OpenWeather API")
	ErrServiceUnavailable   = errors.New("OpenWeather API service is unavailable")
	ErrUnexpectedStatusCode = errors.New("unexpected status code from OpenWeather API")
	ErrDecodingResponse     = errors.New("error decoding response from OpenWeather API")
	ErrTimeout              = errors.New("request to OpenWeather API timed out")
)

type WeatherAPI interface {
	FetchWeatherData(ctx context.Context, lat, lon, openWeatherMapAPIURL, unitsOfMeasurement string) (model.WeatherData, error)
}

type weatherAPI struct{}

func NewWeatherAPI() WeatherAPI {
	return &weatherAPI{}
}

// Define a constant for the timeout duration
const requestTimeout = 5 * time.Second

// FetchWeatherData makes an HTTP request to the OpenWeather API to get weather data for a specific latitude and longitude.
func (api *weatherAPI) FetchWeatherData(ctx context.Context, lat, lon, openWeatherMapAPIURL, unitsOfMeasurement string) (model.WeatherData, error) {
	var data model.WeatherData

	apiKey, ok := ctx.Value(middleware.APIKeyContextKey("apiKey")).(string)
	if !ok || apiKey == "" {
		return data, fmt.Errorf("%w: API key not found in context", ErrBadRequest)
	}

	finalURL, err := util.BuildOpenWeatherMapURL(openWeatherMapAPIURL, apiKey, lat, lon, unitsOfMeasurement)
	if err != nil {
		return data, err
	}

	// Create a new context with a timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(timeoutCtx, "GET", finalURL, nil)
	if err != nil {
		return data, fmt.Errorf("%w: %v", ErrBadRequest, err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		// Check if the error is a timeout
		if errors.Is(err, context.DeadlineExceeded) {
			return data, ErrTimeout
		}
		return data, fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		switch response.StatusCode {
		case http.StatusUnauthorized:
			return data, ErrInvalidAPIKey
		case http.StatusBadRequest:
			return data, ErrBadRequest
		case http.StatusServiceUnavailable:
			return data, ErrServiceUnavailable
		default:
			return data, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, response.StatusCode)
		}
	}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("%w: %v", ErrDecodingResponse, err)
	}

	return data, nil
}
