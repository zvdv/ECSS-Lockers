package email

import (
	"github.com/zvdv/ECSS-Lockers/internal"
	"gopkg.in/gomail.v2"
)

var (
	mailDialier *gomail.Dialer
)

func init() {
	mailDialier = gomail.NewDialer(internal.Env.MailServer,
		internal.Env.MailPort,
		internal.Env.HostEmail,
		internal.Env.HostPassword)
}

func Send(messages ...*gomail.Message) error {
	for _, msg := range messages {
		msg.SetHeader("From", internal.Env.HostEmail)
	}
	return mailDialier.DialAndSend(messages...)
}
