package logger

type LoggerConfig struct {
	Level       string `mapstructure:"level" json:"level" envconfig:"LEVEL"`
	Environment string `mapstructure:"env" json:"env" envconfig:"ENVIRONMENT"`
}
