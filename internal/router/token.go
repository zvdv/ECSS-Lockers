package router

import (
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/templates"
)

func (router *App) tokenValidator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data = struct {
		Token string
	}{
		Token: r.URL.Query().Get("token"),
	}

	if err := templates.Html(w, "templates/validate.html", data); err != nil {
		writeResponse(w, http.StatusInternalServerError, []byte(err.Error()))
	}
}

func (router *App) apiTokenValidator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	logger.Trace("token:%s", token)

	w.Header().Add("HX-Redirect", "/dash")

func makeTokenFromEmail(email string) (string, error) {
	// block schema: 8 byte timestamp, seconds in UTC, expires in 15 mins
	// (900 seconds) this check is done on the receiving end.
	// + len(email)
	buf := make([]byte, len(email)+8)
	binary.BigEndian.PutUint64(buf[:8], uint64(time.Now().Unix()))
	copy(buf[8:], email)

	ciphertext, err := crypto.Encrypt(internal.Env.CipherKey, buf, nil)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ciphertext), nil
}

// returns email, token created time, error
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
