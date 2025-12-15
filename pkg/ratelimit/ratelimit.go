package ratelimit

import (
	"time"

	"go.uber.org/ratelimit"
)

// Limiter abstracts the rate limiting functionality
type Limiter interface {
	Allow() bool
}

type RateLimiter struct {
	limiter ratelimit.Limiter
}

// New creates a new rate limiter that allows `rate` events per second.
func New(rate int) Limiter {
	return &RateLimiter{limiter: ratelimit.New(rate)}
}

func NewWithDuration(rate int, duration time.Duration) Limiter {
	return &RateLimiter{limiter: ratelimit.New(rate, ratelimit.Per(duration))}
}

func NewUnlimited() Limiter {
	return &RateLimiter{limiter: ratelimit.NewUnlimited()}
}

func (rl *RateLimiter) Allow() bool {
	t := rl.limiter.Take()
	return time.Now().Before(t)
}
