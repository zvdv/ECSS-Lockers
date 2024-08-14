package main

import (
	"log"
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal/router"
)

const addr string = "127.0.0.1:8080"

func main() {
	app := router.New()
	log.Printf("Listening at http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.Router()))
}
