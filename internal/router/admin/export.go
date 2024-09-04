package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zvdv/ECSS-Lockers/internal/httputil"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/time"
)

func Export(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.WriteResponse(w, http.StatusMethodNotAllowed, nil)
		return
	}

	lockers, err := queryAllRegistrations()
	if err != nil {
		logger.Error.Println(err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	currentTerm := time.GetCurrentTerm().Format("200601")
	value := fmt.Sprintf("attachment; filename=registrations_%s.csv", currentTerm)
	w.Header().Add("Content-Disposition", value)
	httputil.WriteResponse(w, http.StatusOK, toCSV(lockers))
}

func toCSV(lockers []registration) []byte {
	buf := make([]string, len(lockers)+1)
	buf[0] = ",Locker,Name,Email,Expire On, Email Sent"

	for i, locker := range lockers {
		sent := "false"
		if locker.EmailSent {
			sent = "true"
		}
		buf[i+1] = fmt.Sprintf(
			"%d,%s,%s,%s,%s,%s",
			i+1, locker.Locker,
			locker.Name, locker.Email,
			locker.ExpiryTime.Format("2006-01-02 15:04:05 MST"),
			sent)
	}

	return []byte(strings.Join(buf, "\n"))
}
