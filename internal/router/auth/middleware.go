package auth

import (
	"context"
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal/httputil"
)

func AuthenticatedUserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			httputil.WriteTemplatePage(w, nil,
				"templates/auth/session_expired.html", "templates/nav.html")
			return
		}

		ctx := context.WithValue(r.Context(), httputil.SessionID, cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CRSFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
