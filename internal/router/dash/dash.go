package dash

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	stdtime "time"

	"github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/httputil"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/time"
)

type lockerState struct {
	IsAvailable bool
	LockerID    string
}

func Dash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.WriteResponse(w, http.StatusMethodNotAllowed, nil)
		return
	}

	userEmail, err := httputil.ExtractUserEmail(r)
	if err != nil {
		logger.Error("failed to extract user email from token: %v", err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	data := struct {
		HasLocker  bool
		LockerName string
		ExpireAt   string
	}{
		HasLocker:  false,
		LockerName: "",
		ExpireAt:   "",
	}

	db, lock := database.Lock()
	defer lock.Unlock()

	stmt, err := db.Prepare(`
        SELECT locker, expiry
        FROM registration 
        WHERE user = :email 
        LIMIT 1;`)

	if err != nil {
		logger.Fatal(err)
	}

	var expiry stdtime.Time

	err = stmt.QueryRow(sql.Named("email", userEmail)).Scan(&data.LockerName, &expiry)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error querying for registration: %v", err)
			httputil.WriteResponse(w, http.StatusInternalServerError, nil)

			return
		}

	} else {
		data.HasLocker = true
		data.ExpireAt = expiry.Format("Jan 2, 2006 at 3:04pm")
		httputil.WriteTemplatePage(w, data,
			"templates/dash/index.html",
			"templates/nav.html")

		return
	}

	httputil.WriteTemplatePage(w, data,
		"templates/nav.html",
		"templates/dash/index.html")
}

func ApiLocker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		logger.Error("failed to parse form: %v", err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	locker := r.FormValue("locker")
	if len(locker) == 0 {
		httputil.WriteResponse(w, http.StatusOK, nil)
		return
	}

	lockerNum, err := strconv.ParseUint(locker, 10, 16)
	if err != nil {
		httputil.WriteResponse(
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
        LIKE ?;`)

	if err != nil {
		logger.Fatal("stmt error:", err)
	}

	locker = fmt.Sprintf("%%ELW %d%%", lockerNum)

	rows, err := stmt.Query(locker)
	if err != nil {
		panic(err)
	}

	lockers := []lockerState{}
	for rows.Next() {
		var (
			lockerID       string
			registrationID sql.NullString
		)

		if err := rows.Scan(&lockerID, &registrationID); err != nil {
			logger.Error("failed to scan data: %v", err)
			httputil.WriteResponse(w, http.StatusInternalServerError, nil)
			return
		}

		lockers = append(lockers, lockerState{
			IsAvailable: !registrationID.Valid,
			LockerID:    lockerID,
		})
	}

	data := struct {
		LockerOK bool
		Lockers  []lockerState
	}{
		LockerOK: len(lockers) != 0,
		Lockers:  lockers,
	}

	httputil.WriteTemplateComponent(w, data, "templates/dash/locker_card.html")
}

func DashLockerRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodGet {
		httputil.WriteResponse(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		logger.Error("error parsing form: %v", err)
		httputil.WriteResponse(w, http.StatusBadRequest, nil)
		return
	}

	locker := r.FormValue("locker")

	if r.Method == http.MethodGet {
		httputil.WriteTemplatePage(w, locker,
			"templates/nav.html", "templates/dash/locker_register.html")
		return
	}

	userName := r.FormValue("name")

	userEmail, err := httputil.ExtractUserEmail(r)
	if err != nil {
		logger.Error("error decrypting user email: %v", err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	db, lock := database.Lock()
	defer lock.Unlock()

	var stmt *sql.Stmt

	stmt, err = db.Prepare(`
        SELECT COUNT(*) 
        FROM registration 
        WHERE locker = :locker;`)

	if err != nil {
		logger.Fatal(err)
	}

	var registrationCount uint8

	err = stmt.QueryRow(sql.Named("locker", locker)).Scan(&registrationCount)
	if err != nil {
		logger.Error("error querying for locker: %v", err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	if registrationCount != 0 {
		httputil.WriteTemplateComponent(w, nil, "templates/dash/locker_unavailable.html")
		return
	}

	stmt, err = db.Prepare(`
        INSERT INTO registration (locker, user, name, expiry)
        VALUES (:locker, :user, :name, :expiry);`)

	if err != nil {
		logger.Fatal(err)
	}

	expiryDate := time.NextExpiryDate(time.Now())

	_, err = stmt.Exec(
		sql.Named("locker", locker),
		sql.Named("user", userEmail),
		sql.Named("name", userName),
		sql.Named("expiry", expiryDate))

	if err != nil {
		logger.Error("error writing registration to db: %v", err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	httputil.WriteTemplateComponent(w, nil, "templates/dash/locker_register_ok.html")
}
