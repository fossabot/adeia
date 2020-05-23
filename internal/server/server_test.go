package server

import (
	"adeia-api/internal/config"
	log "adeia-api/internal/logger"
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
	config.Set(&config.Config{
		Logger: config.LoggerConfig{Level: "debug"},
	})
	_ = log.InitLogger()
}

func TestNewAPIServer(t *testing.T) {
	want := &APIServer{Srv: httprouter.New()}

	got := NewAPIServer()

	if !reflect.DeepEqual(got, want) {
		t.Error("should return new APIServer")
	}
}
