package internal

import (
	"github.com/joho/godotenv"
	"github.com/zvdv/ECSS-Lockers/internal/env"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

var (
	Domain    string
	CipherKey []byte
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Warn("failed to load env variables: %v", err)
	}

	Domain = env.MustEnv("DOMAIN")
	CipherKey = []byte(env.MustEnv("CIPHER_KEY"))
}
