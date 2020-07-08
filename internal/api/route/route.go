package route

import (
	"net/http"

	"adeia-api/internal/api/middleware"

	"github.com/go-chi/chi"
)

// Route represents an API route containing a path, method, handler and a
// middleware chain.
type Route struct {
	Method     string
	Path       string
	Handler    http.Handler
	Middleware middleware.FuncChain
}

// New creates a new Route with the provided params.
func New(method, path string, handler http.HandlerFunc, middleware middleware.FuncChain) *Route {
	return &Route{
		Method:     method,
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	}
}

// BindRoutes binds/mounts the provided routes to the router and, also composes
// (adds) all the middleware funcs (in order) to the handler.
func BindRoutes(router chi.Router, routes []*Route) {
	for _, route := range routes {
		handler := route.Handler

		// apply middleware on the handle
		handler = route.Middleware.Compose(handler)

		// mount route
		router.Method(route.Method, route.Path, handler)
	}
}
