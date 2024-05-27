package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/victorbrugnolo/golang-rate-limiter/internal/entity"
)

func RateLimiterMiddleware(ctx context.Context, repository entity.RateLimiterDataRepositoryInterface, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientId := ""
		token := r.Header.Get("API_KEY")
		rateLimiterConfig := entity.RateLimiterConfig{}

		if token != "" {
			clientId = token
			config, err := getRateLimiterConfigByToken(token)

			if err != nil {
				buildErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}

			rateLimiterConfig = *config
		} else {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)

			if err != nil {
				http.Error(w, "Error on getting client IP", http.StatusInternalServerError)
				return
			}

			clientId = ip
			rateLimiterConfig = *getRateLimiterConfigByIP()
		}

		rateLimiterHandler := NewRateLimiter(&rateLimiterConfig, repository)

		if ok := rateLimiterHandler.HandleRateLimit(ctx, clientId); !ok {
			buildErrorResponse(w, http.StatusTooManyRequests, "you have reached the maximum number of requests or actions allowed within a certain time frame")
			return
		}

		next(w, r)
	})
}

func buildErrorResponse(w http.ResponseWriter, httpStatus int, messageBody string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	message := entity.ApiResponse{
		Message: messageBody,
	}

	err := json.NewEncoder(w).Encode(message)

	if err != nil {
		http.Error(w, "Error on enconding response", http.StatusUnprocessableEntity)
	}
}

func getRateLimiterConfigByToken(token string) (*entity.RateLimiterConfig, error) {
	envPrefix := "RATE_LIMITER_BY_TOKEN_" + token
	limit, _ := strconv.Atoi(os.Getenv(envPrefix + "_LIMIT"))
	window, _ := strconv.Atoi(os.Getenv(envPrefix + "_WINDOW"))
	blockWindow, _ := strconv.Atoi(os.Getenv(envPrefix + "_BLOCK_WINDOW"))

	if limit == 0 || window == 0 || blockWindow == 0 {
		return nil, fmt.Errorf("rate limiter configuration for token %s not found", token)
	}

	return entity.NewRateLimiterConfig(limit, time.Duration(window)*time.Second, time.Duration(blockWindow)*time.Second), nil
}

func getRateLimiterConfigByIP() *entity.RateLimiterConfig {
	limit, _ := strconv.Atoi(os.Getenv("RATE_LIMITER_BY_IP_LIMIT"))
	window, _ := strconv.Atoi(os.Getenv("RATE_LIMITER_BY_IP_WINDOW"))
	blockWindow, _ := strconv.Atoi(os.Getenv("RATE_LIMITER_BY_IP_BLOCK_WINDOW"))

	return entity.NewRateLimiterConfig(limit, time.Duration(window)*time.Second, time.Duration(blockWindow)*time.Second)
}
