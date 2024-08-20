package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/email"
	"gopkg.in/gomail.v2"
)

func uvicEmailValidator(email string) bool {
	err := validator.New().Var(email, "email")
	if err == nil {
		return strings.HasSuffix(email, "@uvic.ca")
	}
	return false
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		writeResponse(w, http.StatusBadRequest, []byte("invalid form data"))
		return
	}

	// validate userEmail
	// TODO: may not have to do this? since the incoming userEmail is forced
	// to be @uvic.ca
	userEmail := r.FormValue("email")
	userEmail += "@uvic.ca" // append @uvic since the input data is just netlink id
	if !uvicEmailValidator(userEmail) {
		data := `
            <button type="submit" class="btn btn-primary btn-block">Login</button>
            <div class="form-error">Invalid UVic email address</div>
            `
		writeResponse(w, http.StatusOK, []byte(data))
		return
	}

	tok, err := makeTokenFromEmail(userEmail)
	if err != nil {
		panic(err)
	}

	msg := gomail.NewMessage()
	// msg.SetHeader("From", internal.Env.HostEmail)
	msg.SetHeader("To", userEmail)
	msg.SetHeader("Subject", "Locker registration")
	msg.SetBody("text/html", fmt.Sprintf(emailtemplate,
		internal.Domain,
		tok,
		email.HostEmail))

	if err := email.Send(msg); err != nil {
		panic(err)
	}

	// response
	html := fmt.Sprintf(`<span class="form-info">
        Login link sent to %s!
        </span>`, userEmail)
	writeResponse(w, http.StatusOK, []byte(html))
}

const emailtemplate string = `Hello!
<br />
<br />
You recently requested to sign in to Locker Registration. Click the link below to access your account:
<br />
<br />
<a href="%s/token?token=%s">Sign In to Locker</a>
<br />
<br />
This link will expire in 15 minutes. If you did not request this sign-in, please ignore this email.
<br />
If you need any help, our support team is here for you at <a href="mailto:%s">support</a>.
<br />
<br />
Best regards,
<br />
The Locker Team
<br />`
