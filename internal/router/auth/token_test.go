package auth_test

import (
	"crypto/rand"
	"testing"

	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/router/auth"
)

func TestTokenMaker(t *testing.T) {
	var key [32]byte
	if _, err := rand.Reader.Read(key[:]); err != nil {
		panic(err)
	}

	internal.CipherKey = key[:]

	email := "halnguyen@uvic.ca"
	tok, err := auth.MakeTokenFromEmail(email)
	if err != nil {
		t.Fatal(err)
	}

	email_, _, err := auth.ParseToken(tok)
	if err != nil {
		t.Fatal(err)
	}

	if email != email_ {
		t.Fatal("wrong email out")
	}
}
