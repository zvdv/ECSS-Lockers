package auth

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/zvdv/ECSS-Lockers/internal/crypto"
	"github.com/zvdv/ECSS-Lockers/internal/httputil"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

var (
	AuthRouteRegex *regexp.Regexp = regexp.MustCompile(`/auth/`)
)

func AuthenticatedUserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(string(httputil.SessionID))
		if err != nil {
			w.Header().Add("HX-Redirect", "/sessionexpired")
			return
		}

		if err := cookie.Valid(); err != nil {
			w.Header().Add("HX-Redirect", "/sessionexpired")
			return
		}

		ctx := context.WithValue(r.Context(), httputil.SessionID, cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CRSFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dstURL := r.URL.String()

		if strings.EqualFold(dstURL, "/") ||
			strings.Contains(dstURL, "/auth/") {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		// Check for message cookie
		// get digest
		ok, err := signatureOK(r)
		if err != nil {
			logger.Error.Println(err)
			w.Header().Add("HX-Redirect", "/sessionexpired")
			return
		}

		if !ok {
			logger.Warn.Println("invalid signature")
			w.Header().Add("HX-Redirect", "/sessionexpired")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func signatureOK(r *http.Request) (bool, error) {
	digest, err := getBytesFromCookie(r, string(httputil.Token))
	if err != nil {
		return false, err
	}

	session, err := getBytesFromCookie(r, string(httputil.SessionID))
	if err != nil {
		return false, err
	}

	email, err := crypto.Decrypt(crypto.CipherKey[:], session, nil)
	if err != nil {
		return false, err
	}

	return crypto.VerifySignature(crypto.SignatureKey[:], email, digest)
}

func getBytesFromCookie(r *http.Request, cookieKey string) ([]byte, error) {
	cookie, err := r.Cookie(cookieKey)
	if err != nil {
		return nil, err
	}

	digestB64, err := cookie.Value, cookie.Valid()
	if err != nil {
		return nil, err
	}

	return crypto.Base64.DecodeString(digestB64)
}
