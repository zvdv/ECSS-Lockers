package router

import (
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/crypto"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/templates"
)

func tokenValidator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data = struct {
		Token string
	}{
		Token: r.URL.Query().Get("token"),
	}

	w.WriteHeader(http.StatusOK)
	if err := templates.Html(w, "templates/validate.html", data); err != nil {
		writeResponse(w, http.StatusInternalServerError, []byte(err.Error()))
	}
}

func apiTokenValidator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")

	email, ts, err := parseToken(token)
	if err != nil {
		logger.Error("Failed to parse token:\n%v", err)
		writeResponse(w, http.StatusInternalServerError, nil)
		return
	}

	// token expired?
	now := uint64(time.Now().Unix())
	if now-ts >= 900 { // token expires in 15 mins
		writeResponse(w, http.StatusUnauthorized, []byte("token exired"))
		return
	}

	// cipher email, that will be the auth token
	cookieValue, err := crypto.Encrypt(internal.Env.CipherKey, []byte(email), nil)
	if err != nil {
		panic(err)
	}

	cookie := http.Cookie{
		Name:   "session",
		Value:  crypto.Base64.EncodeToString(cookieValue),
		Domain: internal.Env.Domain,
		Path:   "/",
		MaxAge: 3600, // good for 1 hour
		// MaxAge:   999999, // for local dev
		Secure:   false, // TODO: flip to true on prod
		HttpOnly: false, // TODO: flip to true on prod
	}

	http.SetCookie(w, &cookie)

	w.Header().Add("HX-Redirect", "/dash")
}

// AUTH TOKEN CONSTRUCTION
// block schema: 8 byte timestamp, seconds in UTC, expires in 15 mins
// (900 seconds) this check is done on the receiving end.
// + len(email)
// since the token is transfered via emails as url, hex encode
// instead of base64

// returns hex-string encoded of a token
func makeTokenFromEmail(email string) (string, error) {
	buf := make([]byte, len(email)+8)
	binary.BigEndian.PutUint64(buf[:8], uint64(time.Now().Unix()))
	copy(buf[8:], email)

	ciphertext, err := crypto.Encrypt(internal.Env.CipherKey, buf, nil)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ciphertext), nil
}

// returns email, token created time, parsed from hex-string encoded `token`
func parseToken(token string) (string, uint64, error) {
	decodedTokenBytes, err := hex.DecodeString(token)
	if err != nil {
		return "", 0, err
	}

	pt, err := crypto.Decrypt(internal.Env.CipherKey, decodedTokenBytes, nil)
	if err != nil {
		return "", 0, err
	}

	return string(pt[8:]), binary.BigEndian.Uint64(pt[:8]), nil
}
