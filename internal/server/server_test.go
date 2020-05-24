package server

import (
	log "adeia-api/internal/logger"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/magiconair/properties/assert"
	config "github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	initLogger()
	code := m.Run()
	os.Exit(code)
}

func initLogger() {
	config.Set("logger.level", "debug")
	_ = log.InitLogger()
}

func TestNewAPIServer(t *testing.T) {
	want := &APIServer{Srv: httprouter.New()}

	got := NewAPIServer()

	assert.Equal(t, got, want, "should return new APIServer")
}
