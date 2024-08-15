package main

import (
	"net/http"

	_ "github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/router"
)

const addr string = "127.0.0.1:8080"

func main() {
	app := router.New()
	logger.Trace("Listening at http://%s", addr)
	logger.Info("for local dev, use http://127.0.0.1:8080, for more information, see: https://stackoverflow.com/a/1188145/19114163")
	logger.Fatal(http.ListenAndServe(addr, app.Router()))
}
