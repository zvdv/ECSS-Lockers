package main

import (
	"fmt"
	"net/http"

	_ "github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/database"
	_ "github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/env"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/router"
)

const addr string = "127.0.0.1:8080"

func init() {
	dbURL := fmt.Sprintf(
		"%s?authToken=%s",
		env.MustEnv("DATABASE_URL"),
		env.MustEnv("DATABASE_AUTH_TOKEN"))
	database.Connect(dbURL)
}

func main() {
	mux := router.New()
	logger.Info("Listening at http://%s", addr)
	logger.Info("for local dev, use http://127.0.0.1:8080, for more information, see: https://stackoverflow.com/a/1188145/19114163")
	logger.Fatal(http.ListenAndServe(addr, mux))
}
