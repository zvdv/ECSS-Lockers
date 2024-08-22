package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"hash"

	"github.com/zvdv/ECSS-Lockers/internal/env"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"golang.org/x/crypto/chacha20poly1305"
	"lukechampine.com/blake3"
)

var (
	Base64 = base64.RawStdEncoding

	CipherKey [32]byte
	HMACKey   [32]byte
)

func getKey(key string, buf []byte) {
	if len(buf) != 32 {
		panic("invalid buf length")
	}

	val, err := Base64.DecodeString(env.MustEnv(key))
	if err != nil {
		panic(err)
	}

	if len(val) != 32 {
		panic("invalid key length")
	}

	copy(buf, val)
}

func Initialize() {
	cipherkey, err := Base64.DecodeString(env.MustEnv("CIPHER_KEY"))
	if err != nil {
		logger.Fatal("error decoding cipherkey:", err)
	}
	if len(cipherkey) != 32 {
		logger.Fatal("invalid key length")
	}

	getKey("HMAC_KEY", HMACKey[:])
	getKey("CIPHER_KEY", CipherKey[:])
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

func makeHash() hash.Hash {
	// 32 is the default digest size
	return blake3.New(32, nil)
}

// if you want to append the digest to the end of the message
// pass, call:
//
// `Sign(key, message, message)`
func SignHMAC(key, message, buf []byte) ([]byte, error) {
	var err error

	hash := hmac.New(makeHash, key)

	_, err = hash.Write(message)
	if err != nil {
		return nil, err
	}

	return hash.Sum(buf), nil
}

func VerifyHMAC(key, message, mac []byte) (bool, error) {
	hash := hmac.New(makeHash, key)

	if _, err := hash.Write(message); err != nil {
		return false, err
	}

	expected := hash.Sum(nil)

	return hmac.Equal(expected, mac), nil
}
