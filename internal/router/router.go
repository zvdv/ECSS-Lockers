package router

import (
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/router/auth"
)

func New() *chi.Mux {
	app := chi.NewRouter()

	app.Use(middleware.RealIP)
	app.Use(requestLogger)
	app.Use(middleware.Recoverer)

	app.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	app.Handle("/", http.HandlerFunc(index))

	app.Route("/auth", func(r chi.Router) {
		r.Handle("/", http.HandlerFunc(auth.Auth))
		r.Handle("/api/login", http.HandlerFunc(auth.AuthApiLogin))
		r.Handle("/api/token", http.HandlerFunc(auth.AuthApiToken))
	})

	app.Route("/dash", func(r chi.Router) {
		r.Use(auth.AuthenticatedUserOnly)
		r.Handle("/", http.HandlerFunc(dash))
		r.Handle("/api/locker", http.HandlerFunc(apiLocker))
		r.Handle("/api/locker/confirm", http.HandlerFunc(apiLockerConfirm))
	})

	return app
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		next.ServeHTTP(ww, r)
		end := time.Since(start)

		url := r.URL.Path

		statusCode := ww.Status()
		statusString := color.GreenString("%d", statusCode)
		if statusCode >= 400 {
			statusString = color.RedString("%d", statusCode)
		}

		logger.Trace("%s %s %s from %s - %s %dB in %v",
			r.Method,
			url,
			r.Proto,
			r.RemoteAddr,
			statusString,
			ww.BytesWritten(),
			end)
	})
}
