package router

import (
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/templates"
)

type dashData struct {
	Terms []string
}

func (router *App) dash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data := dashData{
		Terms: []string{"202409", "202501"},
	}

	if err := templates.Html(w, "templates/dash.html", data); err != nil {
		logger.Error(err.Error())
		writeResponse(w, http.StatusInternalServerError, nil)
	}
}

func (router *App) apiDashTerm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	term := r.URL.Query().Get("term")
	writeResponse(w, http.StatusOK, []byte(term))
}
