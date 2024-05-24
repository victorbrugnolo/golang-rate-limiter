package web

import (
	"sync"
	"time"
)

type RateLimiter struct {
	Requests map[string]int
	Limit    int
	Window   time.Duration
	Mutex    sync.Mutex
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		Requests: make(map[string]int),
		Limit:    limit,
		Window:   window,
		Mutex:    sync.Mutex{},
	}
}

func (rl *RateLimiter) RateLimit(clientID string) bool {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()

	if count, exists := rl.Requests[clientID]; !exists || count < rl.Limit {
		if !exists {
			go rl.ResetCount(clientID)
		}

		rl.Requests[clientID]++
		return true
	}

	return false
}

func (rl *RateLimiter) ResetCount(clientID string) {
	time.Sleep(rl.Window)
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()

	delete(rl.Requests, clientID)
}
