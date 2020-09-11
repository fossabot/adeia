package log

import (
	"testing"

	"adeia-api/internal/config"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseLevel(t *testing.T) {
	t.Run("should return correct level when string is valid", func(t *testing.T) {
		want := zap.InfoLevel

		got, err := parseLevel("info")
		assert.Nil(t, err, "should not return error when string is valid")
		assert.Equal(t, want, got, "should return correct level")
	})

	t.Run("should return error when string is invalid", func(t *testing.T) {
		_, err := parseLevel("info123")
		assert.Error(t, err, "should return error when string is invalid")
	})
}

func TestInit(t *testing.T) {
	t.Run("should not return any error when config is valid", func(t *testing.T) {
		c := config.LoggerConfig{
			Level: "debug",
		}

		err := Init(&c)
		assert.Nil(t, err, "should not return any error when config is valid")
	})

	t.Run("should return error when level is invalid", func(t *testing.T) {
		c := config.LoggerConfig{
			Level: "debug123",
		}

		err := Init(&c)
		assert.Error(t, err, "should return error when level is invalid")
	})

	t.Run("should return error when config is invalid", func(t *testing.T) {
		c := config.LoggerConfig{
			Level: "debug",
			Paths: []string{""},
		}

		err := Init(&c)
		assert.Error(t, err, "should return error when config is invalid")
	})
}

func TestSet(t *testing.T) {
	want := zap.NewExample().Sugar()

	logger = &Logger{}
	Set(want)
	got := logger.SugaredLogger

	assert.Equal(t, want, got, "should return the logger set using `set`")
}
