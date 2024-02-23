package main

import (
	"log"
	"net/http"

	"github.com/golang2go/demo-app/weather-service-api/internal/config"
	"github.com/golang2go/demo-app/weather-service-api/internal/handler"
	"github.com/golang2go/demo-app/weather-service-api/internal/middleware"
	"github.com/golang2go/demo-app/weather-service-api/internal/repo"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.DefaultConfig()
	weatherAPI := repo.NewWeatherAPI()
	weatherHandler := handler.NewWeatherHandler(weatherAPI, cfg)

	router := mux.NewRouter()
	router.Use(middleware.RateLimitMiddleware(cfg.RateLimitPerSecond))
	router.Use(middleware.OpenWeatherMapAuthMiddleware)
	router.Use(middleware.LoggingMiddleware)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/weather", weatherHandler.GetWeatherConditionByCoordinates).Methods("GET")

	log.Printf("Starting server on port %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
