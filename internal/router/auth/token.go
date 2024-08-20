package auth

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/time"
	"github.com/zvdv/ECSS-Lockers/internal/crypto"
)

// AUTH TOKEN CONSTRUCTION
// block schema: 8 byte timestamp, seconds in UTC, expires in 15 mins
// (900 seconds) this check is done on the receiving end.
// + len(email)
// since the token is transfered via emails as url, hex encode
// instead of base64

// returns hex-string encoded of a token
func MakeTokenFromEmail(email string) (string, error) {
	buf := make([]byte, len(email)+8)
	binary.BigEndian.PutUint64(buf[:8], uint64(time.Now().Unix()))
	copy(buf[8:], email)

	ciphertext, err := crypto.Encrypt(internal.CipherKey, buf, nil)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ciphertext), nil
}

// returns email, token created time, parsed from hex-string encoded `token`
func ParseToken(token string) (string, uint64, error) {
	decodedTokenBytes, err := hex.DecodeString(token)
	if err != nil {
		return "", 0, err
	}

	pt, err := crypto.Decrypt(internal.CipherKey, decodedTokenBytes, nil)
	if err != nil {
		return "", 0, err
	}

	return string(pt[8:]), binary.BigEndian.Uint64(pt[:8]), nil
}
