# Weather Service

## Overview

Weather Service is a coding exercise that leverages the Open Weather API to provide weather conditions and temperature assessments (hot, cold, or moderate) based on latitude and longitude coordinates. This application exposes an HTTP endpoint that accepts lat/long coordinates and returns the current weather condition in that area.

## Getting Started

### Requirements

- Go 1.22 or later
- Docker (optional)

### Running the Service

#### Using Go

1. Clone the repository and navigate to the project directory.
2. Run the service:
   ```bash
   go run cmd/server/main.go
   ```

#### Using Docker

1. Build the Docker image:
   ```bash
   docker build -t weather-service .
   ```
2. Run the container:
   ```bash
   docker run -d -p 8080:8080 --name my-weather-service weather-service
   ```

### API Usage

#### Endpoint

`GET /api/v1/weather`

#### Query Parameters

- `lat` - Latitude (e.g., `36.9198`)
- `lon` - Longitude (e.g., `93.9276`)

#### Headers

- `X-API-Key` - Your Open Weather Map API key. This key is required for making API requests.

#### Example Request

```bash
curl -X GET "http://localhost:8080/api/v1/weather?lat=36.9198&lon=93.9276" -H "X-API-Key: YOUR_API_KEY"
```

#### Response

The response will include the current weather condition (e.g., snow, rain) and a temperature category (hot, cold, moderate) based on the provided coordinates.

### Security Note

The `X-API-Key` header is used to pass the Open Weather Map API key securely. This method ensures the key is not exposed in URLs, preventing it from being cached or logged in server access logs.

### Defaults

The application uses the following reasonable defaults:

- **Port**: `8080`
- **Rate Limit Per Second**: `5`
- **Open Weather Map API URL**: `https://api.openweathermap.org/data/2.5/weather`
- **Unit of Measurement**: `imperial`

The default unit of measurement is imperial. The default latitude and longitude (in example request cURL) are set to Monett, MO (`36.9198° N, 93.9276° W`).

### Considerations on Concurrency
While leveraging Go's concurrency features like goroutines and channels could enhance the efficiency of fetching data from the Open Weather Map API, this exercise prioritizes simplicity. Implementing such patterns would undoubtedly make the service more scalable and responsive but also introduce complexity that's beyond the scope of this proof of concept. This decision reflects a balance between functionality and maintainability, acknowledging the potential for future scalability while maintaining the current focus on core features.
