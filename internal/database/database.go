package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"github.com/zvdv/ECSS-Lockers/internal/env"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

var (
	db     *sql.DB
	dbLock *sync.Mutex
)

func init() {
	tursoURL := fmt.Sprintf(
		"%s?authToken=%s",
		env.MustEnv("DATABASE_URL"),
		env.MustEnv("DATABASE_AUTH_TOKEN"))

	var err error
	db, err = sql.Open("libsql", tursoURL)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Connected to database")

	dbLock = new(sync.Mutex)
}

func Lock() (*sql.DB, *sync.Mutex) {
	dbLock.Lock()
	return db, dbLock
}
