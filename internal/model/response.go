package model

type MainInfo struct {
	Temp float64 `json:"temp"` // Temperature (based on unit of measurement param)
}

type WeatherCondition struct {
	Main string `json:"main"` // Main weather condition (e.g., Clear, Clouds, Rain)
}

type WeatherData struct {
	Main    MainInfo           `json:"main"`
	Weather []WeatherCondition `json:"weather"`
}

type WeatherResponse struct {
	WeatherCondition string `json:"weatherCondition"` // Current weather condition
	TempCategory     string `json:"tempCategory"`     // Temperature category (hot, cold, moderate)
}
