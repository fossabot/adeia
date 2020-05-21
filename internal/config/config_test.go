package config

import (
	"io/ioutil"
	"os"
	"reflect"
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
	}()

	t.Run("should load without any errors", func(t *testing.T) {
		want := &Config{}
		want.Server.Port = "1234"
		want.Server.Host = "test"

		got, err := Load(validConfFile.Name())

		if err != nil {
			t.Errorf("should not return error. %q", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("should return error when file is nonexistent", func(t *testing.T) {
		_, err := Load("/tmp/foo")

		if err == nil {
			t.Error("should return error when file does not exist")
		}
	})

	t.Run("should return error when yaml is invalid", func(t *testing.T) {
		_, err := Load(invalidConfFile.Name())

		if err == nil {
			t.Error("should return error when yaml is invalid")
		}
	})
}
