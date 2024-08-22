package database

import (
	"database/sql"
	"sync"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

var (
	db     *sql.DB
	dbLock *sync.Mutex
)

func Connect(dbURL string) {
	var err error

	db, err = sql.Open("libsql", dbURL)
	if err != nil {
		logger.Error.Fatal(err)
	}

	logger.Info.Println("Connected to database")

	dbLock = new(sync.Mutex)
}

func Lock() (*sql.DB, *sync.Mutex) {
	dbLock.Lock()
	return db, dbLock
}
