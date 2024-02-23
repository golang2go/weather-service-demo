package config

// Reasonable Defaults
const (
	DefaultPort               = "8080"
	DefaultRateLimitPerSecond = 5
	DefaultOpenWeatherMapURL  = "https://api.openweathermap.org/data/2.5/weather"
	DefaultUnitsOfMeasurement = "imperial"
)

// AppConfig holds application configuration
type AppConfig struct {
	Port                 string
	RateLimitPerSecond   int
	OpenWeatherMapAPIURL string
	UnitOfMeasurement    string
}

// DefaultConfig creates a new AppConfig with default settings.
func DefaultConfig() *AppConfig {
	return NewAppConfig(DefaultPort, DefaultRateLimitPerSecond, DefaultOpenWeatherMapURL, DefaultUnitsOfMeasurement)
}

//NewAppConfig creates a new AppConfig with provided settings or defaults.
func NewAppConfig(port string, rateLimit int, apiURL, unit string) *AppConfig {
	if port == "" {
		port = DefaultPort
	}
	if rateLimit == 0 {
		rateLimit = DefaultRateLimitPerSecond
	}
	if apiURL == "" {
		apiURL = DefaultOpenWeatherMapURL
	}
	if unit == "" {
		unit = DefaultUnitsOfMeasurement
	}

	return &AppConfig{
		Port:                 port,
		RateLimitPerSecond:   rateLimit,
		OpenWeatherMapAPIURL: apiURL,
		UnitOfMeasurement:    unit,
	}
}
