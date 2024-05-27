package entity

import (
	"time"
)

type RateLimiterConfig struct {
	Limit       int
	Window      time.Duration
	BlockWindow time.Duration
}

func NewRateLimiterConfig(limit int, window time.Duration, blockWindow time.Duration) *RateLimiterConfig {
	return &RateLimiterConfig{
		Limit:       limit,
		Window:      window,
		BlockWindow: blockWindow,
	}
}
