package web

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/victorbrugnolo/golang-rate-limiter/internal/entity"
)

type RateLimiterHandler struct {
	Config                    *entity.RateLimiterConfig
	Mutex                     sync.Mutex
	RateLimiterDataRepository entity.RateLimiterDataRepositoryInterface
}

func NewRateLimiter(config *entity.RateLimiterConfig, rateLimiterDataRepository entity.RateLimiterDataRepositoryInterface) *RateLimiterHandler {
	return &RateLimiterHandler{
		Config:                    config,
		Mutex:                     sync.Mutex{},
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
		if time.Since(rateLimiterData.LastSeen) > rl.Config.BlockWindow {
			rateLimiterData.Blocked = false
			rateLimiterData.LastSeen = time.Now()
			rateLimiterData.Requests = 1

			rl.RateLimiterDataRepository.Save(ctx, clientID, rateLimiterData)

			return true

		}

		return false
	}

	if time.Since(rateLimiterData.LastSeen) > rl.Config.Window {
		rateLimiterData.Requests = 0

		rl.RateLimiterDataRepository.Save(ctx, clientID, rateLimiterData)
	}

	if rateLimiterData.Requests < rl.Config.Limit {
		rateLimiterData.LastSeen = time.Now()
		rateLimiterData.Requests++
		rl.RateLimiterDataRepository.Save(ctx, clientID, rateLimiterData)

		return true
	}

	rateLimiterData.Blocked = true
	rl.RateLimiterDataRepository.Save(ctx, clientID, rateLimiterData)

	return false
}
