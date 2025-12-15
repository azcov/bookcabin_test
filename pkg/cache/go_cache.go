package cache

import (
	"time"

	go_cache "github.com/patrickmn/go-cache"
)

var DEFAULT_CACHE_EXPIRATION int64 = 5 * 60 * 1000  // 5 minutes in milliseconds
var DEFAULT_CLEANUP_INTERVAL int64 = 10 * 60 * 1000 // 10 minutes in milliseconds

type goCache struct {
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	cache             *go_cache.Cache
}

func NewGoCache(cfg CacheConfig) Cache {
	if cfg.ExpirationMinute <= 0 {
		cfg.ExpirationMinute = DEFAULT_CACHE_EXPIRATION
	}
	if cfg.CleanupIntervalMinute <= 0 {
		cfg.CleanupIntervalMinute = DEFAULT_CLEANUP_INTERVAL
	}
	defaultExp := time.Duration(cfg.ExpirationMinute) * time.Minute
	cleanupInt := time.Duration(cfg.CleanupIntervalMinute) * time.Minute

	c := go_cache.New(defaultExp, cleanupInt)
	return &goCache{
		defaultExpiration: defaultExp,
		cleanupInterval:   cleanupInt,
		cache:             c,
	}
}

// Implement Cache interface methods here
func (gc *goCache) Get(key string) (any, error) {
	// Implementation goes here
	v, found := gc.cache.Get(key)
	if !found {
		return nil, ErrCacheNotFound
	}
	return v, nil
}

func (gc *goCache) Set(key string, value any) error {
	// Implementation goes here
	gc.cache.Set(key, value, gc.defaultExpiration)
	return nil
}

// SetWithExpiration sets a value in the cache with a specific expiration time
// 0 means default expiration
// -1 means no expiration
func (gc *goCache) SetWithExpiration(key string, value any, exp time.Duration) error {
	// Implementation goes here
	gc.cache.Set(key, value, exp)
	return nil
}

func (gc *goCache) Delete(key string) error {
	// Implementation goes here
	gc.cache.Delete(key)
	return nil
}
