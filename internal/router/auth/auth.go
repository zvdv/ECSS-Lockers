package auth

import (
	"fmt"
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/crypto"
	"github.com/zvdv/ECSS-Lockers/internal/email"
	"github.com/zvdv/ECSS-Lockers/internal/httputil"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/time"
	"gopkg.in/gomail.v2"
)

const (
	tokenExpireLimit    uint64 = 900
	sessionCookieMaxAge int    = 3600
)

func AuthApiLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, []byte("invalid form data"))
		return
	}

	// validate userEmail
	// append @uvic since the input data is just netlink id
	userEmail := r.FormValue("email") + "@uvic.ca"
	if !email.ValidUVicEmail(userEmail) {
		data := `
            <button type="submit" class="btn btn-primary btn-block">Login</button>
            <div class="form-error">Invalid UVic email address</div>
            `
		httputil.WriteResponse(w, http.StatusOK, []byte(data))
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
	httputil.WriteResponse(w, http.StatusOK, []byte(html))
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

	httputil.WriteTemplatePage(w, data,
		"templates/auth/validate.html",
		"templates/nav.html")
}

func AuthApiToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")

	email, ts, err := ParseToken(token)
	if err != nil {
		logger.Error.Printf("Failed to parse token:\n%v\n", err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	// token expired?
	now := uint64(time.Now().Unix())
	if now-ts >= tokenExpireLimit { // token expires in 15 mins
		httputil.WriteResponse(w, http.StatusOK, []byte("token exired"))
		return
	}

	// cipher email, that will be the auth token
	cookieValue, err := crypto.Encrypt(crypto.CipherKey[:], []byte(email), nil)
	if err != nil {
		logger.Error.Printf("error encrypting email: %v\n", err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	cookie := http.Cookie{
		Name:     string(httputil.SessionID),
		Value:    crypto.Base64.EncodeToString(cookieValue),
		Domain:   internal.Domain,
		Path:     "/",
		MaxAge:   sessionCookieMaxAge, // good for 1 hour
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	// sign hmac with plaintext email
	digest, err := crypto.SignMessage(crypto.SignatureKey[:], []byte(email), nil)
	if err != nil {
		logger.Error.Printf("error signing token: %v\n", err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	cookie = http.Cookie{
		Name:     "token",
		Value:    crypto.Base64.EncodeToString(digest),
		Path:     "/",
		Domain:   internal.Domain,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	w.Header().Add("HX-Redirect", "/dash")
}
