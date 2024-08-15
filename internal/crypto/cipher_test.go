package crypto

import (
	"crypto/rand"
	"testing"
)

func makeBuffer(n int) []byte {
	payload := make([]byte, n, n)
	if _, err := rand.Read(payload); err != nil {
		panic(err)
	}
	return payload
}

func TestEncryptDecrypt(t *testing.T) {
	payload := makeBuffer(128)
	key := makeBuffer(32)
	aad := makeBuffer(128)

	ciphertext, err := Encrypt(key, payload, aad)
	if err != nil {
		t.Fatal(err)
	}

	plaintext, err := Decrypt(key, ciphertext, aad)
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
