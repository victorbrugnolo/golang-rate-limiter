package web

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/victorbrugnolo/golang-rate-limiter/internal/entity"
)

type RateLimiterHandler struct {
	Limit                     int
	Window                    time.Duration
	Mutex                     sync.Mutex
	BlockWindow               time.Duration
	RateLimiterDataRepository entity.RateLimiterDataRepositoryInterface
}

func NewRateLimiter(limit int, window time.Duration, blockWindow time.Duration, rateLimiterDataRepository entity.RateLimiterDataRepositoryInterface) *RateLimiterHandler {
	return &RateLimiterHandler{
		Limit:                     limit,
		Window:                    window,
		Mutex:                     sync.Mutex{},
		BlockWindow:               blockWindow,
		RateLimiterDataRepository: rateLimiterDataRepository,
	}
}

func (rl *RateLimiterHandler) HandleRateLimit(ctx context.Context, clientID string) bool {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()

	rateLimiterData, err := rl.RateLimiterDataRepository.Find(ctx, clientID)

	if err != nil {
		log.Printf("Error on finding rate limiter data: %s\n", err)
	}

	if rateLimiterData.Blocked {
		if time.Since(rateLimiterData.LastSeen) > rl.BlockWindow {
			rateLimiterData.Blocked = false
			rateLimiterData.LastSeen = time.Now()
			rateLimiterData.Requests = 1

			rl.RateLimiterDataRepository.Save(ctx, clientID, rateLimiterData)

			return true

		}

		return false
	}

	if time.Since(rateLimiterData.LastSeen) > rl.Window {
		rateLimiterData.Requests = 0

		rl.RateLimiterDataRepository.Save(ctx, clientID, rateLimiterData)
	}

	if rateLimiterData.Requests < rl.Limit {
		rateLimiterData.LastSeen = time.Now()
		rateLimiterData.Requests++
		rl.RateLimiterDataRepository.Save(ctx, clientID, rateLimiterData)

		return true
	}

	rateLimiterData.Blocked = true
	rl.RateLimiterDataRepository.Save(ctx, clientID, rateLimiterData)

	return false
}
