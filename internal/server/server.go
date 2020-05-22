package server

import (
	"adeia-api/internal/config"
	"adeia-api/internal/controllers"
	"adeia-api/internal/logger"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// APIServer is the struct that holds all of the components that need to be
// injected.
type APIServer struct {
	Srv    *httprouter.Router
	Config *config.Config
}

// NewAPIServer returns a new APIServer with the passed-in config.
func NewAPIServer(c *config.Config) *APIServer {
	logger.Debug("initializing new APIServer")
	return &APIServer{Srv: httprouter.New(), Config: c}
}

// AddRoutes registers the handles to the router.
// TODO: properly fix this
func (a *APIServer) AddRoutes() {
	logger.Debug("registering handles to router")
	a.Srv.GET("/", controllers.Index)
}

// Serve starts the server on the host and port, specified in the config.
func (a *APIServer) Serve() error {
	addr := a.Config.Server.Host + ":" + a.Config.Server.Port
	logger.Infof("starting server on %q", addr)
	return http.ListenAndServe(addr, a.Srv)
}
