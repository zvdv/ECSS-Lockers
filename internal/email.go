package internal

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendMail(email string, token string) error {
	d := gomail.NewDialer(Env.MailServer, Env.MailPort, Env.HostEmail, Env.HostPassword)
	msg := gomail.NewMessage()
	msg.SetHeader("From", Env.HostEmail)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Locker registration")
	msg.SetBody("text/html", fmt.Sprintf(`
Hello!
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
<br />
        `, Env.Domain, token, Env.HostEmail))

	return d.DialAndSend(msg)
}
