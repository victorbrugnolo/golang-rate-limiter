package web

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	Requests    map[string]int
	Limit       int
	Window      time.Duration
	Mutex       sync.Mutex
	LastSeen    time.Time
	BlockWindow time.Duration
	Blocked     bool
}

func NewRateLimiter(limit int, window time.Duration, blockWindow time.Duration) *RateLimiter {
	return &RateLimiter{
		Requests:    make(map[string]int),
		Limit:       limit,
		Window:      window,
		Mutex:       sync.Mutex{},
		LastSeen:    time.Now(),
		BlockWindow: blockWindow,
		Blocked:     false,
	}
}

func (rl *RateLimiter) RateLimit(clientID string) bool {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()

	if rl.Blocked {
		if time.Since(rl.LastSeen) > rl.BlockWindow {
			rl.Blocked = false
			delete(rl.Requests, clientID)

			rl.Requests[clientID]++

			return true

		}
		return false
	}

	if count, exists := rl.Requests[clientID]; !exists || count < rl.Limit {
		rl.LastSeen = time.Now()

		// if time.Since(rl.LastSeen) > rl.Window {
		// 	delete(rl.Requests, clientID)
		// }

		rl.Requests[clientID]++
		return true
	}

	return false
}

func (rl *RateLimiter) ResetCount(clientID string) {
	for {
		time.Sleep(rl.Window)

		rl.Mutex.Lock()

		if rl.Blocked {
			fmt.Println(time.Since(rl.LastSeen))
			fmt.Println(rl.BlockWindow)

			if time.Since(rl.LastSeen) > rl.BlockWindow {
				rl.Blocked = false
				delete(rl.Requests, clientID)

			}
		} else {
			delete(rl.Requests, clientID)
		}

		rl.Mutex.Unlock()
	}
}
