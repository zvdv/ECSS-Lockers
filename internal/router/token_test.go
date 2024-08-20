package router

import (
	"testing"

	"github.com/zvdv/ECSS-Lockers/internal"
)

const testCipherKey string = "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"

func TestTokenMaker(t *testing.T) {
	if len(testCipherKey) != 32 {
		panic(len(testCipherKey))
	}
	internal.CipherKey = []byte(testCipherKey)

	email := "halnguyen@uvic.ca"
	tok, err := makeTokenFromEmail(email)
	if err != nil {
		t.Fatal(err)
	}

	email_, _, err := parseToken(tok)
	if err != nil {
		t.Fatal(err)
	}

	if email != email_ {
		t.Fail()
	}
}
