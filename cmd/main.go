package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/victorbrugnolo/golang-rate-limiter/internal/entity"
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

func rateLimiterMiddleware(rateLimiter *web.RateLimiter, next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)

		// go rateLimiter.ResetCount(ip)

		if err != nil {
			http.Error(w, "Error on getting client IP", http.StatusInternalServerError)
			return
		}

		if ok := rateLimiter.RateLimit(ip); !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			message := entity.ApiResponse{
				Message: "you have reached the maximum number of requests or actions allowed within a certain time frame",
			}

			rateLimiter.Blocked = true

			err := json.NewEncoder(w).Encode(message)

			if err != nil {
				http.Error(w, "Error on enconding response", http.StatusUnprocessableEntity)
			}

			return
		}

		next(w, r)
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %s", string(err.Error()))
	}

	limit, _ := strconv.Atoi(os.Getenv("RATE_LIMITER_BY_IP_LIMIT"))
	window, _ := strconv.Atoi(os.Getenv("RATE_LIMITER_BY_IP_WINDOW"))
	blockWindow, _ := strconv.Atoi(os.Getenv("RATE_LIMITER_BY_IP_BLOCK_WINDOW"))

	windowDuration := time.Duration(window) * time.Second
	blockWindowDuration := time.Duration(blockWindow) * time.Second

	rateLimiter := web.NewRateLimiter(limit, windowDuration, blockWindowDuration)

	http.HandleFunc("/ping", rateLimiterMiddleware(rateLimiter, endpointHandler))
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
