package main

import (
	"adeia-api/internal/server"
	"fmt"
	"os"
)

func main() {
	apiServer := server.NewAPIServer()

	apiServer.AddRoutes()

	if err := apiServer.Serve(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%q\n", err)
		os.Exit(1)
	}
}
