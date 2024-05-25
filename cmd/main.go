package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/victorbrugnolo/golang-rate-limiter/internal/entity"
	"github.com/victorbrugnolo/golang-rate-limiter/internal/infra/database"
	"github.com/victorbrugnolo/golang-rate-limiter/internal/infra/web"
)

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	message := entity.ApiResponse{
		Message: "Request successful",
	}

	err := json.NewEncoder(w).Encode(message)

	if err != nil {
		http.Error(w, "Error on enconding response", http.StatusUnprocessableEntity)
	}
}

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %s", string(err.Error()))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		DB:       0,
		Password: "",
	})

	repository := database.NewRedisRateLimiterDataRepository(rdb)

	http.HandleFunc("/ping", web.RateLimiterMiddleware(ctx, repository, endpointHandler))
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
