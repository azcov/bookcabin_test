package ratelimit

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// Limiter abstracts the rate limiting functionality
type Limiter interface {
	Allow() bool
}

type RateLimiter struct {
	limiter *rate.Limiter
}

// New creates a new rate limiter that allows `rate` events per second.
func New(ratePerSecond int) Limiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(
			rate.Limit(ratePerSecond),
			ratePerSecond, // burst = rate
		),
	}
}

// NewWithDuration allows `rate` events per `duration`
func NewWithDuration(rateCount int, duration time.Duration) Limiter {
	r := rate.Every(duration / time.Duration(rateCount))

	return &RateLimiter{
		limiter: rate.NewLimiter(
			r,
			rateCount, // burst
		),
	}
}

// NewUnlimited returns a limiter that never blocks
func NewUnlimited() Limiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Inf, 0),
	}
}

func (rl *RateLimiter) Allow() bool {
	return rl.limiter.Allow()
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	return rl.limiter.Wait(ctx)
}
