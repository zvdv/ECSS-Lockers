package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/zvdv/ECSS-Lockers/internal"
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

	tok, err := makeTokenFromEmail(email)
	if err != nil {
		panic(err)
	}

	if err := internal.SendMail(email, tok); err != nil {
		panic(err)
	}

	// response
	html := fmt.Sprintf(`<span class="form-info">
        Login link sent to %s!
        </span>`, email)
	writeResponse(w, http.StatusOK, []byte(html))
}
