package internal

import (
	"github.com/zvdv/ECSS-Lockers/internal/crypto"
	"github.com/zvdv/ECSS-Lockers/internal/env"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

var (
	Domain    string
	CipherKey []byte
)

func Initialize() {
	Domain = env.Env("DOMAIN")

	cipherKeyString := env.Env("CIPHER_KEY")

	var err error
	CipherKey, err = crypto.Base64.DecodeString(cipherKeyString)

	if err != nil {
		logger.Fatal(err)
	}

	if len(CipherKey) != 32 {
		logger.Fatal("invalid key length.")
	}
}
