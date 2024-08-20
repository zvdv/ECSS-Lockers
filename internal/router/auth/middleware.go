package auth

import (
	"context"
	"net/http"

	"github.com/zvdv/ECSS-Lockers/templates"
)

func AuthenticatedUserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			templates.Html(w, "templates/auth/session_expired.html", nil)
			return
		}

		ctx := context.WithValue(r.Context(), "session_id", cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
