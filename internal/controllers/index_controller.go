package controllers

import (
	"adeia-api/internal/middleware"
	"adeia-api/internal/route"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func IndexRoutes() []*route.Route {
	commonChain := middleware.NewChain(middleware.Logger)

	routes := []*route.Route{
		route.New(http.MethodGet, "/", Index, commonChain),
		route.New(http.MethodGet, "/test", Index2, middleware.Nil),
	}
	return routes
}

// Index is a simple handler that writes a welcome message.
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprint(w, "Welcome\n")
}

// Index2 is a simple handler that writes a welcome message.
func Index2(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprint(w, "Welcome2\n")
}
