package configz

import (
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

var (
	once sync.Once
)

func LoadConfig[T any](cfg *T) error {
	var err error
	once.Do(func() {
		err = envconfig.Process("", &cfg)
		if err != nil {
			log.Fatalf("Failed to load config: %+v", err.Error())
		}
	})
	return err
}
