package config

import (
	"sync"

	"github.com/azcov/bookcabin_test/pkg/cache"
	"github.com/azcov/bookcabin_test/pkg/httpz"
	"github.com/azcov/bookcabin_test/pkg/logger"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Http   httpz.HttpConfig    `mapstructure:"http" json:"http" env:"HTTP"`
	Cache  cache.CacheConfig   `mapstructure:"cache" json:"cache" env:"CACHE"`
	Logger logger.LoggerConfig `mapstructure:"logger" json:"logger" env:"LOGGER"`
}

func NewConfig() *Config {
	return &Config{
		Http: httpz.HttpConfig{
			Port: 8081,
		},
		Cache: cache.CacheConfig{
			Enabled:               true,
			ExpirationMinute:      5,
			CleanupIntervalMinute: 10,
		},
		Logger: logger.LoggerConfig{
			Level:       "info",
			Environment: "development",
		},
	}
}

var (
	once sync.Once
)

func LoadConfig(c *Config) error {
	var err error
	once.Do(func() {
		err = envconfig.Process("", c)
		if err != nil {
			logger.Fatal("Failed to load config: ", "err", err.Error())
		}
	})
	return err
}
