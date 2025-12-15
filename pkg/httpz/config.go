package httpz

type HttpConfig struct {
	Port int `mapstructure:"port" json:"port" envconfig:"PORT"`
}
