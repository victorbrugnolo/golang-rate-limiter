package entity

import "time"

type RateLimiterData struct {
	Requests int
	LastSeen time.Time
	Blocked  bool
}
