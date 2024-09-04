package router

import (
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal/httputil"
)

func Home(w http.ResponseWriter, r *http.Request) {
	// Cache for 15 mins
	w.Header().Add("Cache-Control", "max-age=900")
	w.WriteHeader(http.StatusOK)

	// Parse the template files
	httputil.WriteTemplatePage(w, nil, "templates/index.html")
}

func SessionExpired(w http.ResponseWriter, r *http.Request) {
	// Cache for 15 mins
	w.Header().Add("Cache-Control", "max-age=900")
	w.WriteHeader(http.StatusOK)

	// Parse the template files
	httputil.WriteTemplatePage(w, nil, "templates/auth/session_expired.html", "templates/nav.html")
}
