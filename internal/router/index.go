package router

import (
	"html/template"
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/templates"
)

func (router *App) index(w http.ResponseWriter, r *http.Request) {
	loginTmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		panic(err)
	}

	if err := templates.Html(w, loginTmpl, nil); err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
