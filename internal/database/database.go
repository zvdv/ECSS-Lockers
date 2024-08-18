package database

import (
	"database/sql"
	"sync"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"github.com/zvdv/ECSS-Lockers/internal"
)

var (
	db     *sql.DB
	dbLock *sync.Mutex
)

func Initialize() error {
	dbLock = new(sync.Mutex)
	var err error
	db, err = sql.Open("libsql", internal.Env.DatabaseURL)
	return err
}

func Lock() (*sql.DB, *sync.Mutex) {
    dbLock.Lock()
	return db, dbLock
}
