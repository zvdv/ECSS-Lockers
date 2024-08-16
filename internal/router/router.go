package router

import (
	"net/http"
	"time"

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

	app.router.Use(middleware.RealIP)
	app.router.Use(requestLogger)
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
	if writeData != nil {
		if _, err := w.Write(writeData); err != nil {
			logger.Error("failed to write response: %s", err)
		}
	}
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		next.ServeHTTP(ww, r)
		end := time.Since(start)

		logger.Trace("%s %s %s from %s - %d %dB in %v",
			r.Method,
			r.URL.String(),
			r.Proto,
			r.RemoteAddr,
			ww.Status(),
			ww.BytesWritten(),
			end)
	})
}
