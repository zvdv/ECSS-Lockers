package httputil

import (
	"net/http"

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
