package entity

import "context"

type RateLimiterDataRepositoryInterface interface {
	Find(ctx context.Context, key string) (*RateLimiterData, error)
	Save(ctx context.Context, key string, data *RateLimiterData) error
}
