package httputil

import (
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/crypto"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

func WriteResponse(w http.ResponseWriter, status int, writeData []byte) {
	w.WriteHeader(status)
	if writeData != nil {
		if _, err := w.Write(writeData); err != nil {
			logger.Error("failed to write response: %s", err)
		}
	}
}

func ExtractUserID(r *http.Request) string {
	sessionID, ok := r.Context().Value("session_id").(string)
	if !ok {
		logger.Fatal("ExtractUserToken called from an unprotected route")
	}

	return sessionID
}

func ExtractUserEmail(r *http.Request) (string, error) {
	session := ExtractUserID(r)

	sessionID, err := crypto.Base64.DecodeString(session)
	if err != nil {
		return "", err
	}

	email, err := crypto.Decrypt(internal.CipherKey, sessionID, nil)
	if err != nil {
		return "", err
	}

	return string(email), nil
}
