package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/chacha20poly1305"
)

var (
	Base64 = base64.RawStdEncoding
)

// key length 32 bytes
func Encrypt(key, payload, aad []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	// encryption buffer's length:
	// len(nonce) + len(payload) + len(auth tag)
	nonceSize := aead.NonceSize()
	buf := make([]byte, nonceSize, nonceSize+len(payload)+aead.Overhead())

	// random bytes for nonce
	if _, err := rand.Read(buf[:nonceSize]); err != nil {
		return nil, err
	}

	return aead.Seal(buf, buf, payload, aad), nil
}

// key length 32 bytes
func Decrypt(key, payload, aad []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	nonceSize := aead.NonceSize()
	nonce := payload[:nonceSize]
	ciphertext := payload[nonceSize:]

	return aead.Open(nil, nonce, ciphertext, aad)
}
