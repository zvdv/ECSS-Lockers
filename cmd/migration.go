package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

const schemas string = `
CREATE TABLE IF NOT EXISTS locker (
    id varchar(255) NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS registration (
    locker varchar(255) NOT NULL,
    user varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    expiry datetime NOT NULL,
    expiryEmailSent datetime DEFAULT NULL,
    PRIMARY KEY (locker)
);

CREATE INDEX IF NOT EXISTS user_registration 
ON registration (user);
`

func main() {
	logger.Info("DATABASE MIGRATION")

	dbUrl := fmt.Sprintf("%s?authToken=%s",
		internal.EnvMust("TURSO_DATABASE_URL"),
		internal.EnvMust("TURSO_AUTH_TOKEN"))

	if err := database.Initialize(dbUrl); err != nil {
		logger.Fatal(err)
	}

	db, lock := database.Lock()
	defer lock.Unlock()
	logger.Info("Connected to database")

	if _, err := db.Exec(schemas); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Schema migrated.")

	for i := 1; i <= 200; i++ {
		locker := fmt.Sprintf("ELW %03d", i)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := db.ExecContext(
			ctx,
			"INSERT INTO locker (id) VALUES (:id);",
			sql.Named("id", locker))
		if err != nil {
			logger.Error(
				"Failed to insert locker %s to database:\n%v",
				locker,
				err)
		}
	}
	logger.Info("Seeded lockers")

}
