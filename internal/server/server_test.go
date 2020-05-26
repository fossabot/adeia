package server

import (
	"adeia-api/internal/middleware"
	log "adeia-api/internal/utils/logger"
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
	want := &APIServer{Srv: httprouter.New(), GlobalMiddleware: middleware.NewChain()}

	got := NewAPIServer()

	assert.Equal(t, want, got, "should return new APIServer")
}
