package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/crypto"
	"github.com/zvdv/ECSS-Lockers/internal/email"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/router/utils"
	"github.com/zvdv/ECSS-Lockers/templates"
	"gopkg.in/gomail.v2"
)

func AuthApiLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, []byte("invalid form data"))
		return
	}

	// validate userEmail
	// TODO: may not have to do this? since the incoming userEmail is forced
	// to be @uvic.ca
	userEmail := r.FormValue("email")
	userEmail += "@uvic.ca" // append @uvic since the input data is just netlink id
	if !email.ValidUVicEmail(userEmail) {
		data := `
            <button type="submit" class="btn btn-primary btn-block">Login</button>
            <div class="form-error">Invalid UVic email address</div>
            `
		utils.WriteResponse(w, http.StatusOK, []byte(data))
		return
	}

	// make email token
	tok, err := MakeTokenFromEmail(userEmail)
	if err != nil {
		panic(err)
	}

	msg := gomail.NewMessage()
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
	utils.WriteResponse(w, http.StatusOK, []byte(html))
}

const emailtemplate string = `Hello!
<br />
<br />
You recently requested to sign in to Locker Registration. Click the link below to access your account:
<br />
<br />
<a href="%sauth?token=%s">Sign In to Locker</a>
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

func Auth(w http.ResponseWriter, r *http.Request) {
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
	if err := templates.Html(w, "templates/auth/validate.html", data); err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, []byte(err.Error()))
	}
}

func AuthApiToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")

	email, ts, err := ParseToken(token)
	if err != nil {
		logger.Error("Failed to parse token:\n%v", err)
		utils.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	// token expired?
	now := uint64(time.Now().Unix())
	if now-ts >= 900 { // token expires in 15 mins
		utils.WriteResponse(w, http.StatusUnauthorized, []byte("token exired"))
		return
	}

	// cipher email, that will be the auth token
	cookieValue, err := crypto.Encrypt(internal.CipherKey, []byte(email), nil)
	if err != nil {
		panic(err)
	}

	cookie := http.Cookie{
		Name:   "session",
		Value:  crypto.Base64.EncodeToString(cookieValue),
		Domain: internal.Domain,
		Path:   "/",
		MaxAge: 3600, // good for 1 hour
		// MaxAge:   999999, // for local dev
		Secure:   false, // TODO: flip to true on prod
		HttpOnly: false, // TODO: flip to true on prod
	}

	http.SetCookie(w, &cookie)

	w.Header().Add("HX-Redirect", "/dash")
}
