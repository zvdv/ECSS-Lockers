package router

import (
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/templates"
)

func (router *App) tokenValidator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data = struct {
		Token string
	}{
		Token: r.URL.Query().Get("token"),
	}

	if err := templates.Html(w, "templates/validate.html", data); err != nil {
		writeResponse(w, http.StatusInternalServerError, []byte(err.Error()))
	}
}

func (router *App) apiTokenValidator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	logger.Trace("token:%s", token)

	w.Header().Add("HX-Redirect", "/dash")
}
