package main

import (
	"os"

	"adeia-api/internal/api"
	"adeia-api/internal/config"
	"adeia-api/internal/util/constants"

	_ "github.com/jackc/pgx/v4/stdlib" // Postgres driver
)

func main() {
	// load config
	confPath := getEnv(constants.EnvConfPathKey, "config/config.yaml")
	conf, err := config.Load(confPath)
	checkErr(err)

	// start API server
	checkErr(api.Start(conf))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// getEnv returns value from env if key is present, otherwise returns fallback.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
