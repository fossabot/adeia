package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

const testValidConfPath = "/tmp/adeia-test-valid-config"
const testInvalidConfPath = "/tmp/adeia-test-invalid-config"

func setupTestConf() {
	validConf := `
server:
  host: "test"
  port: 1234
`
	invalidConf := `
@
`
	_ = ioutil.WriteFile(testValidConfPath, []byte(validConf), 0644)
	_ = ioutil.WriteFile(testInvalidConfPath, []byte(invalidConf), 0644)
}

func TestLoad(t *testing.T) {
	setupTestConf()
	defer func() {
		_ = os.Remove(testValidConfPath)
		_ = os.Remove(testInvalidConfPath)
	}()

	t.Run("should load without any errors", func(t *testing.T) {
		want := &Config{}
		want.Server.Port = "1234"
		want.Server.Host = "test"

		got, err := Load(testValidConfPath)

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
		_, err := Load(testInvalidConfPath)

		if err == nil {
			t.Error("should return error when yaml is invalid")
		}
	})
}
