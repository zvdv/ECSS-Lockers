package auth

import (
	"context"
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/crypto"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/templates"
)

func AuthenticatedUserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			templates.Html(w, "templates/auth/session_expired.html", nil)
			return
		}

		sessionID, err := crypto.Base64.DecodeString(cookie.Value)
		if err != nil {
			logger.Error(cookie.Value)
			logger.Fatal("invalid session id:", err)
		}

		email, err := crypto.Decrypt(internal.CipherKey, sessionID, nil)
		if err != nil {
			logger.Fatal("invalid decryption:", err)
		}

		ctx := context.WithValue(r.Context(), "user_email", string(email))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
