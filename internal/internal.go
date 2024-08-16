package internal

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

var (
	Env struct {
		HostEmail    string
		HostPassword string
		MailServer   string
		MailPort     int
		Domain       string
		CipherKey    []byte
	}
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Warn("failed to load .env")
	}

	Env.HostEmail = EnvOrPanic("GMAIL_USER")
	Env.HostPassword = EnvOrPanic("GMAIL_PASSWORD")
	Env.MailServer = "smtp.gmail.com"
	Env.MailPort = 587
	Env.Domain = EnvOrPanic("ORIGIN")
	Env.CipherKey = []byte(EnvOrPanic("CIPHER_KEY"))
	if len(Env.CipherKey) != 32 {
		logger.Warn("invalid value set for env $CIPHER_KEY, expected length of 32 bytes, got %d byte(s)",
			len(Env.CipherKey))
	}
}

func EnvOrPanic(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		logger.Warn("env variable not set: $%s", key)
	}
	return value
}
