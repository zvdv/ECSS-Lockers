package router

import (
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/templates"
)

func index(w http.ResponseWriter, r *http.Request) {
	// cache for 15 mins
	w.Header().Add("Cache-Control", "max-age=900")
	w.WriteHeader(http.StatusOK)
	if err := templates.Html(w, "templates/index.html", nil); err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
