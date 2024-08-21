package email

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/zvdv/ECSS-Lockers/internal/env"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"gopkg.in/gomail.v2"
)

var (
	mailDialier *gomail.Dialer
	HostEmail   string
)

func Initialize() {
	HostEmail = env.Env("EMAIL_HOST_ADDRESS")

	mailDialier = gomail.NewDialer("smtp.gmail.com",
		587,
		HostEmail,
		env.Env("EMAIL_HOST_PASSWORD"))
}

// NOTE: overides the "From" header, will set it to
// $EMAIL_HOST_ADDRESS email address.
func Send(messages ...*gomail.Message) error {
	for _, msg := range messages {
		msg.SetHeader("From", HostEmail)
	}
	return mailDialier.DialAndSend(messages...)
}

func ValidUVicEmail(email string) bool {
	err := validator.New().Var(email, "email")
	if err == nil {
		return strings.HasSuffix(email, "@uvic.ca")
	}
	logger.Error("Invalid email %s:\n%v", email, err)
	return false
}
