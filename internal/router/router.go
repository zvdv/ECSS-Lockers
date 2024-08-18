package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

func New() *chi.Mux {
	app := chi.NewRouter()

	app.Use(middleware.RealIP)
	app.Use(requestLogger)
	app.Use(middleware.Recoverer)

	app.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	app.Handle("/", http.HandlerFunc(index))
	app.Handle("/api/login", http.HandlerFunc(login))
	app.Handle("/token", http.HandlerFunc(tokenValidator))
	app.Handle("/api/token", http.HandlerFunc(apiTokenValidator))

	// TODO: Middleware to validate cookie here
	app.Handle("/dash", http.HandlerFunc(dash))
	app.Handle("/api/dash/term", http.HandlerFunc(apiDashTerm))

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
