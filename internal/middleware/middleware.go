package middleware

import (
	log "adeia-api/internal/utils/logger"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Func func(h httprouter.Handle) httprouter.Handle

type FuncChain struct {
	funcs []Func
}

var Nil = FuncChain{}

func NewChain(funcs ...Func) FuncChain {
	return FuncChain{append(([]Func)(nil), funcs...)}
}

func (c *FuncChain) Compose(f httprouter.Handle) httprouter.Handle {
	for _, m := range c.funcs {
		f = m(f)
	}
	return f
}

func (c *FuncChain) Append(funcs ...Func) *FuncChain {
	newChain := make([]Func, 0, len(c.funcs)+len(funcs))
	newChain = append(newChain, c.funcs...)
	newChain = append(newChain, funcs...)

	return &FuncChain{newChain}
}

func Logger(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Infof("path: %q", r.URL.Path)
		next(w, r, p)
	}
}
