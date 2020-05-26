package route

import (
	"adeia-api/internal/middleware"

	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Method     string
	Path       string
	Handle     httprouter.Handle
	Middleware middleware.FuncChain
}

func New(method, path string, handle httprouter.Handle, middleware middleware.FuncChain) *Route {
	return &Route{
		Method:     method,
		Path:       path,
		Handle:     handle,
		Middleware: middleware,
	}
}

func BindRoutes(router *httprouter.Router, routes []*Route) {
	for _, route := range routes {
		handle := route.Handle

		// apply middleware on the handle
		handle = route.Middleware.Compose(handle)

		// mount route
		router.Handle(route.Method, route.Path, handle)
	}
}
