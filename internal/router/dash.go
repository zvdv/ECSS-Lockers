package router

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

func apiLocker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	lockerNum, err := strconv.ParseUint(r.FormValue("locker"), 10, 16)
	if err != nil {
		writeResponse(
			w,
			http.StatusOK,
			[]byte("<p class=\"text-error text-center\">Invalid locker</p>"))
		return
	}

	db, lock := database.Lock()
	defer lock.Unlock()

	stmt, err := db.Prepare(`
        SELECT locker.id, registration.locker FROM locker
        LEFT JOIN registration ON locker.id = registration.locker
        WHERE locker.id LIKE?;
        `)
	if err != nil {
		logger.Fatal("stmt error", err)
		return
	}

	locker := fmt.Sprintf("%%ELW %d%%", lockerNum)
	rows, err := stmt.Query(locker)
	if err != nil {
		panic(err)
	}

	type LockerState struct {
		IsAvailable bool
		LockerID    string
	}

	lockers := []LockerState{}
	for rows.Next() {
		var (
			lockerID       string
			registrationID string
		)

		rows.Scan(&lockerID, &registrationID)
		lockers = append(lockers, LockerState{
			IsAvailable: len(registrationID) == 0,
			LockerID:    lockerID,
		})
	}

	data := struct {
		LockerOK bool
		Lockers  []LockerState
	}{
		LockerOK: len(lockers) != 0,
		Lockers:  lockers,
	}

	if err := templates.Component(w, "templates/dash_locker_card.html", data); err != nil {
		panic(err)
	}
}

func apiLockerConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeResponse(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	locker := r.FormValue("locker")
	logger.Info(locker)

    // TODO: calculate exp timestamp

	// TODO: write to db
}
