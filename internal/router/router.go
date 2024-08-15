package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

type App struct {
	router *chi.Mux
	// TODO: App state
}

func (app *App) Router() *chi.Mux {
	return app.router
}

func New() *App {
	app := &App{chi.NewRouter()}

	app.router.Use(middleware.Logger)
	app.router.Use(middleware.Recoverer)

	app.router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	app.router.Handle("/", http.HandlerFunc(app.index))
	app.router.Handle("/api/login", http.HandlerFunc(app.login))
	app.router.Handle("/token", http.HandlerFunc(app.tokenValidator))
	app.router.Handle("/api/token", http.HandlerFunc(app.apiTokenValidator))

	// TODO: middleware to validate cookie here
	app.router.Handle("/dash", http.HandlerFunc(app.dash))
	app.router.Handle("/api/dash/term", http.HandlerFunc(app.apiDashTerm))

	return app
}

func writeResponse(w http.ResponseWriter, status int, writeData []byte) {
	w.WriteHeader(status)
	if _, err := w.Write(writeData); err != nil {
        logger.Error("failed to write response: %s", err)
	}
}
