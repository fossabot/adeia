package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"adeia-api/internal/util"

	"github.com/spf13/viper"
)

var (
	// initConf is used to ensure that config is initialized only once.
	initConf = new(sync.Once)

	// envOverrides holds all environment value keys for overriding the config.
	envOverrides = map[string]string{
		// mailer overrides
		"mailer.username": util.EnvMailerUsername,
		"mailer.password": util.EnvMailerPassword,

		// database env overrides
		"database.dbname":   util.EnvDBNameKey,
		"database.user":     util.EnvDBUserKey,
		"database.password": util.EnvDBPasswordKey,
		"database.host":     util.EnvDBHostKey,
		"database.port":     util.EnvDBPortKey,

		// cache env overrides
		"cache.host": util.EnvCacheHostKey,
		"cache.port": util.EnvCachePortKey,
	}
)

func setEnvOverrides() {
	viper.SetEnvPrefix(util.EnvPrefix)
	for k, v := range envOverrides {
		// The only error that is returned from this method is when len(input) == 0.
		// So we can safely ignore it.
		_ = viper.BindEnv(k, v)
	}
}

// Load loads YAML from file specified by EnvConfPathKey into viper.
// The file must be a readable file containing valid YAML.
func Load() error {
	err := errors.New("config already loaded")

	initConf.Do(func() {
		err = nil

		confPath := getEnv(util.EnvConfPathKey, "config/config.yaml")
		basePath := filepath.Base(confPath)

		viper.SetConfigName(strings.TrimSuffix(basePath, filepath.Ext(basePath)))
		viper.AddConfigPath(filepath.Dir(confPath))
		viper.SetConfigType("yaml")

		// set env overrides for secrets
		setEnvOverrides()

		e := viper.ReadInConfig()
		if e != nil {
			err = fmt.Errorf("cannot read config: %v", e)
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
