package server

import (
	"adeia-api/internal/controllers"
	"adeia-api/internal/route"
	log "adeia-api/internal/utils/logger"
	"net/http"

	"github.com/julienschmidt/httprouter"
	config "github.com/spf13/viper"
)

// APIServer is the struct that holds all of the components that need to be
// injected.
type APIServer struct {
	Srv *httprouter.Router
}

// NewAPIServer returns a new APIServer with the passed-in config.
func NewAPIServer() *APIServer {
	log.Debug("initializing new APIServer")
	return &APIServer{Srv: httprouter.New()}
}

// AddRoutes registers the handles to the router.
// TODO: properly fix this
func (a *APIServer) AddRoutes() {
	log.Debug("registering handles to router")
	route.BindRoutes(a.Srv, controllers.IndexRoutes())
}

// Serve starts the server on the host and port, specified in the config.
func (a *APIServer) Serve() error {
	addr := config.GetString("server.host") + ":" + config.GetString("server.port")
	log.Infof("starting server on %q", addr)
	return http.ListenAndServe(addr, a.Srv)
}
