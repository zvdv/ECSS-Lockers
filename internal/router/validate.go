package router

import (
	"html/template"
	"log"
	"net/http"

	"github.com/zvdv/ECSS-Lockers/templates"
)

func (router *App) tokenValidator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/validate.html")
	if err != nil {
		panic(err)
	}

	var data = struct {
		Token string
	}{
		Token: r.URL.Query().Get("token"),
	}

	templates.Base(w, tmpl, data)
}

func (router *App) apiTokenValidator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	log.Println("token:", token)

	w.Header().Add("HX-Redirect", "/dash")
}
