package controller

import (
	"fmt"
	"net/http"

	"adeia-api/internal/middleware"
	"adeia-api/internal/route"
)

// IndexRoutes contains all of the route-info needed, to bind to the router.
// It returns a slice of routes that it is responsible for.
func IndexRoutes() []*route.Route {
	routes := []*route.Route{
		route.New(http.MethodGet, "/", Index, middleware.Nil),
		route.New(http.MethodGet, "/test", Index2, middleware.Nil),
	}
	return routes
}

// Index is a simple handler that writes a welcome message.
func Index(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Welcome\n")
}

// Index2 is a simple handler that writes a welcome message.
func Index2(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Welcome 2\n")
}
