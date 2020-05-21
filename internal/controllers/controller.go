package controllers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Index is simple handler that writes a welcome message.
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprint(w, "Welcome\n")
}
