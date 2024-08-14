package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
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
		panic(err)
	}

	// validate email
	email := r.FormValue("email")
	if !uvicEmailValidator(email) {
		w.Write([]byte(`
            <button type="submit" class="btn btn-primary btn-block">Login</button>
            <div class="form-error">Invalid UVic email address</div>
            `))
		return
	}

	// TODO: send email

	// response
	html := fmt.Sprintf(`<span class="form-info">
        Login link sent to %s!<br/>
        </span>`, email)
	w.Write([]byte(html))
}
