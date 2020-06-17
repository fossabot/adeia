package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"adeia-api/internal/controller"
	"adeia-api/internal/middleware"
	"adeia-api/internal/route"
	log "adeia-api/internal/utils/logger"
	"adeia-api/internal/utils/ratelimiter"

	"github.com/julienschmidt/httprouter"
	config "github.com/spf13/viper"
	"golang.org/x/time/rate"
)

// Server is the struct that holds all of the components that need to be
// injected.
type Server struct {
	Srv              *httprouter.Router
	GlobalMiddleware middleware.FuncChain
}

// New returns a new Server with the passed-in config.
func New() *Server {
	log.Debug("initializing new API server")

	// create a new ratelimiter
	l := ratelimiter.New(
		rate.Limit(config.GetFloat64("server.ratelimit_rate")),
		config.GetInt("server.ratelimit_burst"),
		time.Duration(config.GetInt("server.ratelimit_window"))*time.Second,
	)

	return &Server{
		Srv:              httprouter.New(),
		GlobalMiddleware: middleware.NewChain(middleware.RateLimiter(l)),
	}
}

// AddRoutes registers the handles to the router.
func (s *Server) AddRoutes() {
	log.Debug("registering handles to router")
	route.BindRoutes(s.Srv, controller.IndexRoutes())
}

// Serve starts the server on the host and port, specified in the config.
func (s *Server) Serve() {
	addr := config.GetString("server.host") + ":" + config.GetString("server.port")
	srv := &http.Server{
		// TODO: add timeouts
		Addr:    addr,
		Handler: s.GlobalMiddleware.Compose(s.Srv),
	}

	go func() {
		log.Infof("starting server on %q", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("error while serving: %v", err)
		}
	}()

	// graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	// wait for 5 seconds for pending requests to be completed, otherwise cancel them
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		// requests that exceed the context deadline (5 seconds) are cancelled
		log.Errorf("failed to gracefully shutdown server: %v", err)
	} else {
		log.Info("server gracefully stopped")
	}
}
