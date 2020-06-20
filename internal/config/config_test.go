package config

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"adeia-api/internal/util"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupTestConf(content, pattern string) (*os.File, error) {
	f, err := ioutil.TempFile("", pattern)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = f.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	return f, nil
}

func TestLoad(t *testing.T) {
	// setup dummy config files in /tmp/
	validConf := `
server:
  host: test
  port: 1234
logger:
  level: info
`
	validConfFile, err := setupTestConf(validConf, "adeia-valid-config")
	if err != nil {
		t.Fatal("error setting up valid test config")
	}

	invalidConf := `
@
`
	invalidConfFile, err := setupTestConf(invalidConf, "adeia-invalid-config")
	if err != nil {
		t.Fatal("error setting up invalid test config")
	}

	// cleanup when test ends
	defer func() {
		_ = os.Remove(validConfFile.Name())
		_ = os.Remove(invalidConfFile.Name())
		_ = os.Unsetenv(util.EnvConfPathKey)
	}()

	t.Run("should load without any errors", func(t *testing.T) {
		want := viper.New()
		want.Set("server.port", "1234")
		want.Set("server.host", "test")
		want.Set("logger.level", "info")

		_ = os.Setenv(util.EnvConfPathKey, validConfFile.Name())
		err := Load()
		got := viper.GetViper()

		assert.Nil(t, err, "should not return error when config is valid")

		assert.Equal(t, want.GetString("server.port"), got.GetString("server.port"), "should be equal")
		assert.Equal(t, want.GetString("server.host"), got.GetString("server.host"), "should be equal")
		assert.Equal(t, want.GetString("logger.level"), got.GetString("logger.level"), "should be equal")
	})

	initConf = new(sync.Once)

	t.Run("should return error when file is nonexistent", func(t *testing.T) {
		_ = os.Setenv(util.EnvConfPathKey, "/tmp/foo")
		err := Load()
		assert.Error(t, err, "should return error when file does not exist")
	})

	initConf = new(sync.Once)

	t.Run("should return error when yaml is invalid", func(t *testing.T) {
		_ = os.Setenv(util.EnvConfPathKey, invalidConfFile.Name())
		err := Load()
		assert.Error(t, err, "should return error when yaml is invalid")
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("should return value from env if key is set", func(t *testing.T) {
		_ = os.Setenv("DUMMY_KEY", "foo")

		want := "foo"
		got := getEnv("DUMMY_KEY", "bar")

		assert.Equal(t, want, got, "should return value from env if key is set")
	})

	_ = os.Unsetenv("DUMMY_KEY")

	t.Run("should return fallback if key is not set", func(t *testing.T) {
		want := "bar"
		got := getEnv("DUMMY_KEY", "bar")

		assert.Equal(t, want, got, "should return fallback if key is not set")
	})
}
