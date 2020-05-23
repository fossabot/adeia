package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

// Config holds all of the config needed for the application.
type Config struct {
	Server ServerConfig `yaml:"server"`
	Logger LoggerConfig `yaml:"logger"`
}

// ServerConfig holds server-specific config.
type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// LoggerConfig holds config for the logger.
type LoggerConfig struct {
	Level string `yaml:"level"`
}

var (
	config   *Config
	initConf *sync.Once
)

func init() {
	initConf = new(sync.Once)
}

// LoadConf loads YAML from confPath into a new Config.
// confPath must be a readable file containing valid YAML.
func LoadConf() error {
	err := errors.New("config already loaded")

	initConf.Do(func() {
		err = nil

		confPath := getEnv("ADEIA_CONF_PATH", "./config/config.yaml")
		conf := &Config{}

		file, e := os.Open(confPath)
		if e != nil {
			err = e
			return
		}
		defer func() {
			cErr := file.Close()
			if err == nil {
				err = cErr
			}
		}()

		d := yaml.NewDecoder(file)
		if e = d.Decode(&conf); e != nil {
			err = e
		}

		config = conf
	})

	return err
}

// Set sets the config.
func Set(c *Config) {
	config = c
}

// Get returns the config instance.
func Get() *Config {
	return config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
