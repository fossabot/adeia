package controllers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Index is a simple handler that writes a welcome message.
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprint(w, "Welcome\n")
}
