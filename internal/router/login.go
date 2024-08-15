package router

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/crypto"
)

func uvicEmailValidator(email string) bool {
	err := validator.New().Var(email, "email")
	if err == nil {
		return strings.HasSuffix(email, "@uvic.ca")
	}
	return false
}

func (router *App) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		writeResponse(w, http.StatusBadRequest, []byte("invalid form data"))
		return
	}

	// validate email
	// TODO: may not have to do this? since the incoming email is forced
	// to be @uvic.ca
	email := r.FormValue("email")
	email += "@uvic.ca" // append @uvic since the input data is just netlink id
	if !uvicEmailValidator(email) {
		data := `
            <button type="submit" class="btn btn-primary btn-block">Login</button>
            <div class="form-error">Invalid UVic email address</div>
            `
		writeResponse(w, http.StatusOK, []byte(data))
		return
	}

	if err := sendLoginLink(email); err != nil {
		panic(err)
	}

	// response
	html := fmt.Sprintf(`<span class="form-info">
        Login link sent to %s!
        </span>`, email)
	writeResponse(w, http.StatusOK, []byte(html))
}

func sendLoginLink(email string) error {
	// block schema: 8 byte timestamp, seconds in UTC, expires in 15 mins
	// (900 seconds) this check is done on the receiving end.
	// + len(email)
	buf := make([]byte, len(email)+8)
	binary.BigEndian.AppendUint64(buf[:8], uint64(time.Now().Unix()))
	copy(buf[8:], email)
	var key [32]byte // TODO: pull this from env CIPHER_KEY

	ciphertext, err := crypto.Encrypt(key[:], buf, nil)
	if err != nil {
		return err
	}

	token := crypto.Base64Encode.EncodeToString(ciphertext)
	return internal.SendMail(email, token)
}
