package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"adeia-api/internal/api/middleware"
	"adeia-api/internal/api/route"
	"adeia-api/internal/cache"
	"adeia-api/internal/controller"
	"adeia-api/internal/db"
	"adeia-api/internal/util/log"
	"adeia-api/internal/util/ratelimiter"

	"github.com/julienschmidt/httprouter"
	config "github.com/spf13/viper"
)

// Server is the struct that holds all of the components that need to be
// injected.
type Server struct {
	cache            cache.Cache
	db               db.DB
	globalMiddleware middleware.FuncChain
	srv              *httprouter.Router
}

// New returns a new Server with the passed-in config.
func New(d db.DB, c cache.Cache) *Server {
	log.Debug("initializing new API server")

	l := getGlobalRateLimiter()
	return &Server{
		srv:              httprouter.New(),
		globalMiddleware: middleware.NewChain(middleware.RateLimiter(l)),
		db:               d,
		cache:            c,
	}
}

// AddRoutes registers the handles to the router.
func (s *Server) AddRoutes() {
	log.Debug("registering handles to router")

	controller.Init(s.db, s.cache)
	route.BindRoutes(s.srv, controller.UserRoutes())
}

// Serve starts the server on the host and port, specified in the config.
func (s *Server) Serve() {
	addr := config.GetString("server.host") + ":" + config.GetString("server.port")
	srv := &http.Server{
		// TODO: add timeouts
		// TODO: add TLS support
		Addr:    addr,
		Handler: s.globalMiddleware.Compose(s.srv),
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

func getGlobalRateLimiter() ratelimiter.RateLimiter {
	// create a new ratelimiter
	// when r = b = x, then user is allowed to make a max. of x requests per second
	r := config.GetInt("server.ratelimit_rate")
	l := ratelimiter.New(
		float64(r),
		r,
		time.Duration(config.GetInt("server.ratelimit_window"))*time.Second,
	)

	return l
}
