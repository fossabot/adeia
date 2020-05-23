package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"sync"
	"testing"
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

	return f, nil
}

func TestLoad(t *testing.T) {
	validConf := `
server:
  host: "test"
  port: 1234
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
		_ = os.Unsetenv("ADEIA_CONF_PATH")
	}()

	t.Run("should load without any errors", func(t *testing.T) {
		want := &Config{}
		want.Server.Port = "1234"
		want.Server.Host = "test"

		_ = os.Setenv("ADEIA_CONF_PATH", validConfFile.Name())
		err := LoadConf()
		got := Get()

		if err != nil {
			t.Errorf("should not return error. %q", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	initConf = new(sync.Once)

	t.Run("should return error when file is nonexistent", func(t *testing.T) {
		_ = os.Setenv("ADEIA_CONF_PATH", "/tmp/foo")
		err := LoadConf()

		if err == nil {
			t.Error("should return error when file does not exist")
		}
	})

	initConf = new(sync.Once)

	t.Run("should return error when yaml is invalid", func(t *testing.T) {
		_ = os.Setenv("ADEIA_CONF_PATH", invalidConfFile.Name())
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
