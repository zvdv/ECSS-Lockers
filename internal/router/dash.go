package router

import (
	"html/template"
	"net/http"

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

	tmpl, err := template.ParseFiles("templates/dash.html")
	if err != nil {
		panic(err)
	}

	data := dashData{
		Terms: []string{"202409", "202501"},
	}

	templates.Html(w, tmpl, data)
}

func (router *App) apiDashTerm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	term := r.URL.Query().Get("term")
	w.Write([]byte(term))
}
