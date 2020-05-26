package server

import (
	"adeia-api/internal/controllers"
	"adeia-api/internal/middleware"
	"adeia-api/internal/route"
	log "adeia-api/internal/utils/logger"
	"net/http"

	"github.com/julienschmidt/httprouter"
	config "github.com/spf13/viper"
)

// APIServer is the struct that holds all of the components that need to be
// injected.
type APIServer struct {
	Srv              *httprouter.Router
	GlobalMiddleware middleware.FuncChain
}

// NewAPIServer returns a new APIServer with the passed-in config.
func NewAPIServer() *APIServer {
	log.Debug("initializing new API server")
	return &APIServer{Srv: httprouter.New(), GlobalMiddleware: middleware.NewChain()}
}

// AddRoutes registers the handles to the router.
func (a *APIServer) AddRoutes() {
	log.Debug("registering handles to router")
	route.BindRoutes(a.Srv, controllers.IndexRoutes())
}

// Serve starts the server on the host and port, specified in the config.
func (a *APIServer) Serve() error {
	addr := config.GetString("server.host") + ":" + config.GetString("server.port")
	log.Infof("starting server on %q", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: a.GlobalMiddleware.Compose(a.Srv),
	}
	return srv.ListenAndServe()
}
