package logger

import (
	"reflect"
	"sync"
	"testing"

	config "github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestParseLevel(t *testing.T) {
	t.Run("should return proper level when string is valid", func(t *testing.T) {
		want := zap.InfoLevel

		got, err := parseLevel("info")
		if err != nil {
			t.Errorf("should not return any error: %v", err)
		}
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("should return error when string is invalid", func(t *testing.T) {
		_, err := parseLevel("info123")
		if err == nil {
			t.Error("should return error when string is invalid")
		}
	})
}

func TestInit(t *testing.T) {
	t.Run("should not return any error when config is valid", func(t *testing.T) {
		config.Set("logger.level", "debug")

		err := InitLogger()
		if err != nil {
			t.Errorf("should not return any error when config is valid: %v", err)
		}
	})

	// reset sync.Once because logger was already initialized in the previous test
	initLog = new(sync.Once)

	t.Run("should return error when config is invalid", func(t *testing.T) {
		config.Set("logger.level", "debug123")

		err := InitLogger()
		if err == nil {
			t.Error("should return error when config is invalid")
		}
	})
}

func TestSet(t *testing.T) {
	want := zap.NewExample().Sugar()

	SetLogger(want)

	got := logger.SugaredLogger

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
