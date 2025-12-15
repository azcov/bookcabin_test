package cache

import "time"

type Cache interface {
	Get(key string) (any, error)
	Set(key string, value any) error
	SetWithExpiration(key string, value any, exp time.Duration) error
	Delete(key string) error
}
