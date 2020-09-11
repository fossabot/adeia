package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"adeia-api/internal/api/middleware"
	"adeia-api/internal/config"
	holidayController "adeia-api/internal/controller/holiday"
	roleController "adeia-api/internal/controller/role"
	userController "adeia-api/internal/controller/user"
	"adeia-api/internal/util/constants"
	"adeia-api/internal/util/log"
	"adeia-api/internal/util/ratelimiter"

	"github.com/go-chi/chi"
)

// Server is the struct that holds all of the components that need to be
// injected.
type Server struct {
	config           *config.ServerConfig
	globalMiddleware middleware.FuncChain
	srv              chi.Router
}

// New returns a new Server with the passed-in config.
func New(conf *config.ServerConfig) *Server {
	log.Debug("initializing new API server")

	l := getGlobalRateLimiter(conf.RateLimitRate, conf.RateLimitWindow)
	return &Server{
		config:           conf,
		globalMiddleware: middleware.NewChain(middleware.RateLimiter(l)),
		srv:              chi.NewRouter(),
	}
}

// AddRoutes registers the handles to the router.
func (s *Server) AddRoutes() {
	log.Debug("registering handles to router")

	// setup router
	s.srv.Route("/"+constants.APIVersion, func(r chi.Router) {
		r.Mount(userController.Routes())
		r.Mount(holidayController.Routes())
		r.Mount(roleController.Routes())
	})
}

// Serve starts the server on the host and port, specified in the config.
func (s *Server) Serve() {
	addr := s.config.Host + ":" + strconv.Itoa(s.config.Port)
	srv := &http.Server{
		// TODO: add timeouts
		// TODO: add TLS support
		Addr:    addr,
		Handler: s.globalMiddleware.Compose(s.srv),
	}

	// make chan for interrupts
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	// make chan for errors
	serverErrors := make(chan error, 1)

	go func() {
		log.Infof("starting server on %q", addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// graceful shutdown
	select {
	case err := <-serverErrors:
		if err != http.ErrServerClosed {
			log.Panicf("error while serving: %v", err)
		}

	case sig := <-sigint:
		log.Infof("main: %v: start shutdown", sig)

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
}

func getGlobalRateLimiter(ratelimitRate, ratelimitWindow int) ratelimiter.RateLimiter {
	// create a new ratelimiter
	// when refillRate = bucketSize, say some value `x`, then user will be allowed to
	// make a max. of `x` requests per second
	limiter := ratelimiter.New(
		float64(ratelimitRate),
		ratelimitRate,
		time.Duration(ratelimitWindow)*time.Second,
	)

	return limiter
}
