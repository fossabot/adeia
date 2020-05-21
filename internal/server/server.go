package server

import (
	"adeia-api/internal/controllers"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type APIServer struct {
	Srv *httprouter.Router
}

func NewAPIServer() *APIServer {
	return &APIServer{httprouter.New()}
}

func (a *APIServer) AddRoutes() {
	a.Srv.GET("/", controllers.Index)
}

func (a *APIServer) Serve() error {
	return http.ListenAndServe(":5000", a.Srv)
}
