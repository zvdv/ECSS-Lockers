package email

import (
	"github.com/zvdv/ECSS-Lockers/internal/env"
	"gopkg.in/gomail.v2"
)

var (
	mailDialier *gomail.Dialer
	HostEmail   string
)

func init() {
	HostEmail = env.MustEnv("EMAIL_HOST_ADDRESS")

	mailDialier = gomail.NewDialer("smtp.gmail.com",
		587,
		HostEmail,
		env.MustEnv("EMAIL_HOST_PASSWORD"))
}

// NOTE: overides the "From" header, will set it to
// $EMAIL_HOST_ADDRESS email address.
func Send(messages ...*gomail.Message) error {
	for _, msg := range messages {
		msg.SetHeader("From", HostEmail)
	}
	return mailDialier.DialAndSend(messages...)
}
