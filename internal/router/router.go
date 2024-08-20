package router

import (
	"context"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/crypto"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/templates"
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
	app.Route("/dash", func(r chi.Router) {
		r.Use(authenticatedUserOnly)
		r.Handle("/", http.HandlerFunc(dash))
		r.Handle("/api/locker", http.HandlerFunc(apiLocker))
		r.Handle("/api/locker/confirm", http.HandlerFunc(apiLockerConfirm))
	})

	return app
}

func authenticatedUserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			err := templates.Html(w, "templates/session_expired.html", nil)
			if err != nil {
				writeResponse(w, http.StatusInternalServerError, nil)
			}
			return
		}

		sessionID, err := hex.DecodeString(cookie.Value)
		if err != nil {
			logger.Fatal("invalid session id:", err)
		}

		email, err := crypto.Decrypt(internal.CipherKey, sessionID, nil)
		if err != nil {
			logger.Fatal("invalid decryption:", err)
		}

		logger.Info("%s", string(email))

		ctx := context.WithValue(r.Context(), "user_email", string(email))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
