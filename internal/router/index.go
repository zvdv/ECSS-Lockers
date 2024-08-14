package router

import (
	"html/template"
	"net/http"

	"github.com/zvdv/ECSS-Lockers/templates"
)

func (router *App) index(w http.ResponseWriter, r *http.Request) {
	loginTmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		panic(err)
	}
	templates.Base(w, loginTmpl, nil)
}
