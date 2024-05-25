package main

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/victorbrugnolo/golang-rate-limiter/cmd/internal/entity"
	"github.com/victorbrugnolo/golang-rate-limiter/cmd/internal/infra/web"
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
	rateLimiter := web.NewRateLimiter(2, 5*time.Second, 10*time.Second)

	http.HandleFunc("/ping", rateLimiterMiddleware(rateLimiter, endpointHandler))
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
