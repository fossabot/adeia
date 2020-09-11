package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestConfigFile(content, pattern string) (*os.File, error) {
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
	validConfFile, err := newTestConfigFile(validConf, "adeia-valid-config")
	if err != nil {
		t.Fatal("error setting up valid test config")
	}

	invalidConf := `
@
`
	invalidConfFile, err := newTestConfigFile(invalidConf, "adeia-invalid-config")
	if err != nil {
		t.Fatal("error setting up invalid test config")
	}

	invalidConf2 := `
server:
  port: hello
`
	invalidConfFile2, err := newTestConfigFile(invalidConf2, "adeia-invalid-config-2")
	if err != nil {
		t.Fatal("error setting up invalid test config")
	}

	// cleanup when tests end
	defer func() {
		_ = os.Remove(validConfFile.Name())
		_ = os.Remove(invalidConfFile.Name())
		_ = os.Remove(invalidConfFile2.Name())
	}()

	t.Run("should load without any errors", func(t *testing.T) {
		want := Config{
			ServerConfig: ServerConfig{
				Port: 1234,
				Host: "test",
			},
			LoggerConfig: LoggerConfig{
				Level: "info",
			},
		}

		got, err := Load(validConfFile.Name())
		assert.Nil(t, err, "should not return error when config is valid")

		assert.Equal(t, want.ServerConfig.Port, got.ServerConfig.Port)
		assert.Equal(t, want.ServerConfig.Host, got.ServerConfig.Host)
		assert.Equal(t, want.LoggerConfig.Level, got.LoggerConfig.Level)
	})

	t.Run("should return error when file is non-existent", func(t *testing.T) {
		_, err := Load("/tmp/foo")
		assert.Error(t, err, "should return error when file does not exist")
	})

	t.Run("should return error when yaml is invalid", func(t *testing.T) {
		_, err := Load(invalidConfFile.Name())
		assert.Error(t, err, "should return error when yaml is invalid")
	})

	t.Run("should return error unmarshalling fails", func(t *testing.T) {
		_, err := Load(invalidConfFile2.Name())
		assert.Error(t, err, "should return error unmarshalling fails")
	})
}
