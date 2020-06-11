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
var initConf = new(sync.Once)

func setEnvOverrides() {
	viper.SetEnvPrefix(utils.EnvPrefix)

	// The only error that is returned from this method is when len(input) == 0.
	// So we can safely ignore them.
	_ = viper.BindEnv("database.dbname", utils.EnvDBNameKey)
	_ = viper.BindEnv("database.user", utils.EnvDBUserKey)
	_ = viper.BindEnv("database.password", utils.EnvDBPasswordKey)
	_ = viper.BindEnv("database.host", utils.EnvDBHostKey)
	_ = viper.BindEnv("database.port", utils.EnvDBPortKey)
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
		setEnvOverrides()

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
