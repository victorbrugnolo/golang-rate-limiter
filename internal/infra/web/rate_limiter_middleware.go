package web

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/victorbrugnolo/golang-rate-limiter/internal/entity"
)

func RateLimiterMiddleware(ctx context.Context, repository entity.RateLimiterDataRepositoryInterface, next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)

		if err != nil {
			http.Error(w, "Error on getting client IP", http.StatusInternalServerError)
			return
		}

		limit, _ := strconv.Atoi(os.Getenv("RATE_LIMITER_BY_IP_LIMIT"))
		window, _ := strconv.Atoi(os.Getenv("RATE_LIMITER_BY_IP_WINDOW"))
		blockWindow, _ := strconv.Atoi(os.Getenv("RATE_LIMITER_BY_IP_BLOCK_WINDOW"))

		windowDuration := time.Duration(window) * time.Second
		blockWindowDuration := time.Duration(blockWindow) * time.Second

		rateLimiterHandler := NewRateLimiter(limit, windowDuration, blockWindowDuration, repository)

		if ok := rateLimiterHandler.HandleRateLimit(ctx, ip); !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			message := entity.ApiResponse{
				Message: "you have reached the maximum number of requests or actions allowed within a certain time frame",
			}

			err = json.NewEncoder(w).Encode(message)

			if err != nil {
				http.Error(w, "Error on enconding response", http.StatusUnprocessableEntity)
			}

			return
		}

		next(w, r)
	})
}
