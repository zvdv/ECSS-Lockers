package admin

import (
	"context"
	"database/sql"
	"net/http"
	stdtime "time"

	"github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/httputil"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

func Registrations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		httputil.WriteResponse(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		logger.Error.Println(err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	locker := r.FormValue("locker")

	db, lock := database.Lock()
	defer lock.Unlock()

	stmt, err := db.Prepare(`DELETE FROM registration WHERE locker = :locker;`)
	if err != nil {
		logger.Error.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*stdtime.Second)
	defer cancel()

	result, err := stmt.ExecContext(ctx, sql.Named("locker", locker))
	if err != nil {
		logger.Error.Println(err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	rowsAff, err := result.RowsAffected()
	logger.Trace.Printf("Deleted locker %s, result: %d row(s), err: %v\n", locker, rowsAff, err)

	status := http.StatusNoContent
	if rowsAff == 0 {
		status = http.StatusNotFound
	}
	httputil.WriteResponse(w, status, nil)
}
