package server

import (
	"adeia-api/internal/config"
	"adeia-api/internal/controllers"
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
	return &APIServer{Srv: httprouter.New(), Config: c}
}

// AddRoutes registers the handles to the router.
// TODO: properly fix this
func (a *APIServer) AddRoutes() {
	a.Srv.GET("/", controllers.Index)
}

// Serve starts the server on the host and port, specified in the config.
func (a *APIServer) Serve() error {
	return http.ListenAndServe(a.Config.Server.Host+":"+a.Config.Server.Port, a.Srv)
}
