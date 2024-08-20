package dash

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/router/ioutil"
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

func Dash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	email, ok := r.Context().Value("user_email").(string)
	if !ok {
		logger.Fatal("credential not found in protected route")
	}

	db, lock := database.Lock()
	defer lock.Unlock()

	data := dashData{
		HasData:  false,
		UserName: "",
		Locker:   "",
		Lockers:  []lockerState{},
	}

	stmt, err := db.Prepare(`
        SELECT locker, name 
        FROM registration 
        WHERE user = :email 
        LIMIT 1;`)
	if err != nil {
		logger.Fatal(err)
	}

	err = stmt.QueryRow(sql.Named("email", email)).Scan(&data.Locker, &data.UserName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		}
	} else {
		data.HasData = true
		templates.Html(w, "templates/dash/index.html", data)
		return
	}

	// query for all lockers
	rows, err := db.Query("SELECT id FROM locker;")
	for rows.Next() {
		var locker string
		rows.Scan(&locker)
		data.Lockers = append(data.Lockers, lockerState{locker, false})
	}

	templates.Html(w, "templates/dash/index.html", data)
}

func ApiLocker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		logger.Error("failed to parse form: %v", err)
		ioutil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	lockerNum, err := strconv.ParseUint(r.FormValue("locker"), 10, 16)
	if err != nil {
		ioutil.WriteResponse(
			w,
			http.StatusOK,
			[]byte("<p class=\"text-error text-center\">Invalid locker</p>"))
		return
	}

	db, lock := database.Lock()
	defer lock.Unlock()

	stmt, err := db.Prepare(`
        SELECT locker.id, registration.locker 
        FROM locker
        LEFT JOIN registration 
        ON locker.id = registration.locker
        WHERE locker.id 
        LIKE ?;
        `)
	if err != nil {
		logger.Fatal("stmt error:", err)
	}

	locker := fmt.Sprintf("%%ELW %d%%", lockerNum)
	rows, err := stmt.Query(locker)
	if err != nil {
		panic(err)
	}

	type Locker struct {
		IsAvailable bool
		LockerID    string
	}

	lockers := []Locker{}
	for rows.Next() {
		var (
			lockerID       string
			registrationID sql.NullString
		)

		if err := rows.Scan(&lockerID, &registrationID); err != nil {
			logger.Error("failed to scan data: %v", err)
			ioutil.WriteResponse(w, http.StatusInternalServerError, nil)
			return
		}

		lockers = append(lockers, Locker{
			IsAvailable: !registrationID.Valid,
			LockerID:    lockerID,
		})
	}

	data := struct {
		LockerOK bool
		Lockers  []Locker
	}{
		LockerOK: len(lockers) != 0,
		Lockers:  lockers,
	}

	templates.Component(w, "templates/dash/locker_card.html", data)
}

func DashLockerRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		ioutil.WriteResponse(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	locker := r.FormValue("locker")

	db, lock := database.Lock()
	defer lock.Unlock()

	var (
		stmt *sql.Stmt
		err  error
	)

	stmt, err = db.Prepare(`
        SELECT COUNT(*) 
        FROM registration 
        WHERE locker = :locker;`)
	if err != nil {
		logger.Fatal(err)
	}

	var registrationCount int
	err = stmt.QueryRow(locker).Scan(&registrationCount)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Trace("%d", registrationCount)

	// TODO: calculate exp timestamp

	// TODO: write to db
}
