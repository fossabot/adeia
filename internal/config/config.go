package config

import (
	"adeia-api/internal/utils"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

// initConf is used to ensure that config is initialized only once.
var initConf *sync.Once

func init() {
	initConf = new(sync.Once)
}

// LoadConf loads YAML from file specified by EnvConfPathKey into viper.
// The file must be a readable file containing valid YAML.
func LoadConf() error {
	err := errors.New("config already loaded")

	initConf.Do(func() {
		err = nil

		confPath := getEnv(utils.EnvConfPathKey, "./config/config.yaml")
		base := filepath.Base(confPath)

		viper.SetConfigName(strings.TrimSuffix(base, filepath.Ext(base)))
		viper.AddConfigPath(filepath.Dir(confPath))
		viper.SetConfigType("yaml")
		viper.SetEnvPrefix("adeia")

		// TODO: add env overrides for cache

		e := viper.ReadInConfig()
		if e != nil {
			err = e
			return
		}
	})

	return err
}

// getEnv returns value from env if key is present, otherwise returns fallback.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
