package server

import (
	"adeia-api/internal/config"
	"adeia-api/internal/logger"
	"os"
	"reflect"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestMain(m *testing.M) {
	initLogger()
	code := m.Run()
	os.Exit(code)
}

func initLogger() {
	_ = logger.Init(&config.LoggerConfig{Level: "debug"})
}

func TestNewAPIServer(t *testing.T) {
	want := &APIServer{Srv: httprouter.New(), Config: &config.Config{}}

	got := NewAPIServer(&config.Config{})

	if !reflect.DeepEqual(got, want) {
		t.Error("should return new APIServer")
	}
}
