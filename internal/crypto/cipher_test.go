package crypto_test

import (
	"crypto/rand"
	"testing"

	"github.com/zvdv/ECSS-Lockers/internal/crypto"
)

func makeBuffer(n int) []byte {
	payload := make([]byte, n)
	if _, err := rand.Read(payload); err != nil {
		panic(err)
	}

	return payload
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	payload := makeBuffer(128)
	key := makeBuffer(32)
	aad := makeBuffer(128)

	ciphertext, err := crypto.Encrypt(key, payload, aad)
	if err != nil {
		t.Fatal(err)
	}

	plaintext, err := crypto.Decrypt(key, ciphertext, aad)
	if err != nil {
		t.Fatal(err)
	}

	if len(payload) != len(plaintext) {
		t.Fatal("invalid length")
	}

	for i := range plaintext {
		if payload[i] != plaintext[i] {
			t.Fatal("ciphertext != plaintext")
		}
	}
}

func TestHMAC(t *testing.T) {
	t.Parallel()

	key := makeBuffer(32)

	message := makeBuffer(128)

	disgest, err := crypto.SignMessage(key, message, nil)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := crypto.VerifySignature(key, message, disgest)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fail()
	}

	disgest[8] ^= 1 // flip a bit
	ok, err = crypto.VerifySignature(key, message, disgest)
	if err != nil {
		t.Fatal(err)
	}

	if ok {
		t.Fail()
	}
}
