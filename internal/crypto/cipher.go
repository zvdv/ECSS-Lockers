package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"

	"github.com/zvdv/ECSS-Lockers/internal/env"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"golang.org/x/crypto/chacha20poly1305"
	"lukechampine.com/blake3"
)

// this needs to be constant and hardcoded, with the format
// [application timestamp purpose].
//
// see Blake3 specs, section 6.2:
// https://raw.githubusercontent.com/BLAKE3-team/BLAKE3-specs/master/blake3.pdf
//
// though this can be changed in the future, since it's used to derive a key for
// CSRF tokens, which are short-lived.
const (
	kdfContextString string = "ECSS Lockers Registration 2024-08-23 11:36:07 Signature"
	SignatureSize    int    = 32
)

var (
	Base64 = base64.RawStdEncoding

	CipherKey    [32]byte
	SignatureKey [32]byte
)

func Initialize() {
	// parse cipher key from env
	envCipherkey, err := Base64.DecodeString(env.MustEnv("CIPHER_KEY"))
	if err != nil {
		logger.Error.Fatal("error decoding cipherkey:", err)
	}
	if len(envCipherkey) != 32 {
		logger.Error.Fatal("invalid key length")
	}

	copy(CipherKey[:], envCipherkey)

	// derive signature key from cipher key
	blake3.DeriveKey(SignatureKey[:], kdfContextString, envCipherkey)
}

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

// if you want to append the digest to the end of the message
// pass, call:
//
// `Sign(key, message, message)`
func SignMessage(key, message, buf []byte) ([]byte, error) {
	var err error

	hash := blake3.New(SignatureSize, SignatureKey[:])
	_, err = hash.Write(message)
	if err != nil {
		return nil, err
	}

	return hash.Sum(buf), nil
}

func VerifySignature(key, message, mac []byte) (bool, error) {
	hash := blake3.New(SignatureSize, SignatureKey[:])

	if _, err := hash.Write(message); err != nil {
		return false, err
	}

	expected := hash.Sum(nil)

	return hmac.Equal(expected, mac), nil
}
