package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"adeia-api/internal/util/constants"

	"github.com/spf13/viper"
)

// Load loads config from confPath into viper. The file must be readable and
// must contain valid YAML.
func Load(confPath string) (*Config, error) {
	v := viper.New()
	basePath := filepath.Base(confPath)

	v.SetConfigName(strings.TrimSuffix(basePath, filepath.Ext(basePath)))
	v.AddConfigPath(filepath.Dir(confPath))
	v.SetConfigType("yaml")

	// set env overrides for secrets
	setEnvOverrides(v, constants.EnvPrefix, envOverrides)

	// read config
	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot read config: %v", err)
	}

	// unmarshal config
	var c Config
	err = v.Unmarshal(&c)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal to config struct: %v", err)
	}

	return &c, nil
}

// setEnvOverrides sets the keys of env variables that can override the config, in viper.
func setEnvOverrides(v *viper.Viper, envPrefix string, overrides map[string]string) {
	v.SetEnvPrefix(envPrefix)
	for key, val := range overrides {
		// The only error that is returned from this method is when len(input) == 0.
		// So we can safely ignore it.
		_ = v.BindEnv(key, val)
	}
}
