package httputil

import (
	"html/template"
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/crypto"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

func WriteTemplateComponent(w http.ResponseWriter, data interface{}, filename ...string) {
	tmpl := template.Must(template.ParseFiles(filename...))

	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error("error executing template data: %v", err)
		WriteResponse(w, http.StatusInternalServerError, nil)
	}
}

func WriteTemplatePage(w http.ResponseWriter, data interface{}, filename ...string) {
	files := make([]string, 1, len(filename)+1)

	files[0] = "templates/base.html"
	files = append(files, filename...)

	tmpl := template.Must(template.ParseFiles(files...))

	w.WriteHeader(http.StatusOK)
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		logger.Error("error executing template data: %v", err)
		WriteResponse(w, http.StatusInternalServerError, nil)
	}
}

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
