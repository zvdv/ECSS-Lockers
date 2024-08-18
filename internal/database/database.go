package database

import (
	"database/sql"
	"sync"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var (
	db     *sql.DB
	dbLock *sync.Mutex
)

func Initialize(dbUrl string) error {
	dbLock = new(sync.Mutex)
	var err error
	db, err = sql.Open("libsql", dbUrl)
	return err
}

func Lock() (*sql.DB, *sync.Mutex) {
	dbLock.Lock()
	return db, dbLock
}
