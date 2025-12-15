package cache

type CacheConfig struct {
	Enabled               bool  `mapstructure:"enabled" json:"enabled" envconfig:"ENABLED"`
	ExpirationMinute      int64 `mapstructure:"expiration_minute" json:"expiration_minute" envconfig:"EXPIRATION_MINUTE"`
	CleanupIntervalMinute int64 `mapstructure:"cleanup_interval_minute" json:"cleanup_interval_minute" envconfig:"CLEANUP_INTERVAL_MINUTE"`
}
