package middleware

import "net/http"

// Func represents a middleware function.
type Func func(h http.Handler) http.Handler

// FuncChain is a slice of Handlers, representing the middleware chain.
type FuncChain struct {
	funcs []Func
}

// Nil represents an empty middleware chain.
var Nil = FuncChain{}

// NewChain creates an empty middleware chain.
func NewChain(funcs ...Func) FuncChain {
	return FuncChain{append(([]Func)(nil), funcs...)}
}

// Compose applies/composes all the middleware funcs in-order on the provided
// handler.
func (c *FuncChain) Compose(f http.Handler) http.Handler {
	for _, m := range c.funcs {
		f = m(f)
	}
	return f
}

// Append appends the passed-in funcs to the existing middleware chain.
func (c *FuncChain) Append(funcs ...Func) FuncChain {
	newChain := make([]Func, 0, len(c.funcs)+len(funcs))
	newChain = append(newChain, c.funcs...)
	newChain = append(newChain, funcs...)

	return FuncChain{newChain}
}
