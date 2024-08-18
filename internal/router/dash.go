package router

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/templates"
)

type dashData struct {
	HasData  bool
	UserName string
	Locker   string
	Lockers  []lockerState
}

type lockerState struct {
	ID    string
	InUse bool
}

func dash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	email, ok := r.Context().Value("user_email").(string)
	if !ok {
		panic("credential not found in protected route")
	}

	logger.Info(email)

	db, lock := database.Lock()
	defer lock.Unlock()

	data := dashData{
		HasData:  false,
		UserName: "",
		Locker:   "",
		Lockers:  []lockerState{},
	}

	err := db.QueryRow(
		`SELECT locker, name FROM registration WHERE user = :email LIMIT 1;`,
		sql.Named("email", email)).Scan(&data.Locker, &data.UserName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		}
	} else {
		data.HasData = true
		if err := templates.Html(w, "templates/dash.html", data); err != nil {
			logger.Error(err.Error())
			writeResponse(w, http.StatusInternalServerError, nil)
		}
		return
	}

	// query for all lockers
	rows, err := db.Query("SELECT id FROM locker;")
	for rows.Next() {
		var locker string
		rows.Scan(&locker)
		data.Lockers = append(data.Lockers, lockerState{locker, false})
	}

	if err := templates.Html(w, "templates/dash.html", data); err != nil {
		logger.Error(err.Error())
		writeResponse(w, http.StatusInternalServerError, nil)
	}
}

func apiDashTerm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	term := r.URL.Query().Get("term")
	writeResponse(w, http.StatusOK, []byte(term))
}
