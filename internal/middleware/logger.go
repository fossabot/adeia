package middleware

import (
	"net/http"

	log "adeia-api/internal/utils/logger"
)

// Logger is a simple middleware that logs the URL.Path of every request.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("path: %q", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Logger2 is a simple middleware that logs the URL.Path of every request.
func Logger2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("path2: %q", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
