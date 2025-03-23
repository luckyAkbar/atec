package usecase

import (
	"context"

	"github.com/go-redis/redis_rate/v10"
)

// RateLimiter interface for rate limiter. at the time of writing, this is only used
// to help unit testing as a mocked implementation of this interface will be easier
// to generate. also to conform to the dependency inversion principle when refined later
type RateLimiter interface {
	Allow(ctx context.Context, key string, limit redis_rate.Limit) (*redis_rate.Result, error)
}
