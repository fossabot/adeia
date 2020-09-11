package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
