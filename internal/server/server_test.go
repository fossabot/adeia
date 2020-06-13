package server

import (
	"os"
	"testing"

	"adeia-api/internal/middleware"
	log "adeia-api/internal/utils/logger"

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
	_ = log.Init()
}

func TestNewAPIServer(t *testing.T) {
	want := &Server{Srv: httprouter.New(), GlobalMiddleware: middleware.NewChain()}

	got := New()

	assert.Equal(t, want, got, "should return new Server")
}
