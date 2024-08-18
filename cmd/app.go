package main

import (
	"fmt"
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal"
	_ "github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/router"
)

const addr string = "127.0.0.1:8080"

func main() {
	tursoUrl := internal.EnvMust("TURSO_DATABASE_URL")
	tursoToken := internal.EnvMust("TURSO_AUTH_TOKEN")

	dbUrl := fmt.Sprintf("%s?authToken=%s", tursoUrl, tursoToken)
	if err := database.Initialize(dbUrl); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connected to database at %s", tursoUrl)

	mux := router.New()
	logger.Info("Listening at http://%s", addr)
	logger.Info("for local dev, use http://127.0.0.1:8080, for more information, see: https://stackoverflow.com/a/1188145/19114163")
	logger.Fatal(http.ListenAndServe(addr, mux))
}
