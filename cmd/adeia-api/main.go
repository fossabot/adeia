package main

import (
	"adeia-api/internal/api"
	"adeia-api/internal/config"
)

func main() {
	// load config
	checkErr(config.Load())

	// start API server
	checkErr(api.Start())
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
