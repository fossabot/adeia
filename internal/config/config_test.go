package config

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/spf13/viper"
)

func setupTestConf(content, pattern string) (*os.File, error) {
	f, err := ioutil.TempFile("", pattern)
	if err != nil {
		return nil, err
	}

	_, err = f.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	_ = f.Chmod(0666)

	return f, nil
}

func TestLoad(t *testing.T) {
	validConf := `
server:
  host: test
  port: 1234

logger:
  level: info
`
	validConfFile, err := setupTestConf(validConf, "adeia-valid-config")
	if err != nil {
		t.Errorf("error setting up test config file")
	}

	invalidConf := `
@
`
	invalidConfFile, err := setupTestConf(invalidConf, "adeia-invalid-config")
	if err != nil {
		t.Errorf("error setting up test config file")
	}

	defer func() {
		_ = os.Remove(validConfFile.Name())
		_ = os.Remove(invalidConfFile.Name())
		_ = os.Unsetenv(EnvConfPathKey)
	}()

	t.Run("should load without any errors", func(t *testing.T) {
		want := viper.New()
		want.Set("server.port", "1234")
		want.Set("server.host", "test")
		want.Set("logger.level", "info")

		t.Log(validConfFile.Name())
		_ = os.Setenv(EnvConfPathKey, validConfFile.Name())
		err := LoadConf()
		got := viper.GetViper()

		if err != nil {
			t.Errorf("should not return error. %v", err)
		}

		if want.GetString("server.port") != got.GetString("server.port") ||
			want.GetString("server.host") != got.GetString("server.host") ||
			want.GetString("logger.level") != got.GetString("logger.level") {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	initConf = new(sync.Once)

	t.Run("should return error when file is nonexistent", func(t *testing.T) {
		_ = os.Setenv(EnvConfPathKey, "/tmp/foo")
		err := LoadConf()

		if err == nil {
			t.Error("should return error when file does not exist")
		}
	})

	initConf = new(sync.Once)

	t.Run("should return error when yaml is invalid", func(t *testing.T) {
		_ = os.Setenv(EnvConfPathKey, invalidConfFile.Name())
		err := LoadConf()

		if err == nil {
			t.Error("should return error when yaml is invalid")
		}
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("should return value from env if key is set", func(t *testing.T) {
		_ = os.Setenv("DUMMY_KEY", "foo")
		want := "foo"

		got := getEnv("DUMMY_KEY", "bar")
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	_ = os.Unsetenv("DUMMY_KEY")

	t.Run("should return fallback if key is not set", func(t *testing.T) {
		want := "bar"
		got := getEnv("DUMMY_KEY", "bar")
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
