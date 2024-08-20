package main

import (
	"database/sql"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/env"
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
    expiryEmailSent boolean DEFAULT FALSE,
    PRIMARY KEY (locker)
);

CREATE INDEX IF NOT EXISTS user_registration 
ON registration (user);
`

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Fatal(err)
	}

	logger.Info("DATABASE MIGRATION")

	database.Connect(fmt.Sprintf(
		"%s?authToken=%s",
		env.MustEnv("DATABASE_URL"),
		env.MustEnv("DATABASE_AUTH_TOKEN")))

	db, lock := database.Lock()
	defer lock.Unlock()

	if _, err := db.Exec(schemas); err != nil {
		logger.Fatal(err)
	}

	logger.Info("created schema.")

	logger.Info("seeding 200 lockers..")
	// eeehhh i'm not proud of how this is being done but
	// database/sql does not support array type for query
	// arg out of the box :(
	for i := 0; i < 200; i++ {
		locker := fmt.Sprintf("ELW %03d", i+1)

		stmt, err := db.Prepare(`INSERT INTO locker (id) VALUES (:id);`)
		if err != nil {
			logger.Fatal(err)
		}

		_, err = stmt.Exec(sql.Named("id", locker))
		if err != nil {
			logger.Error("error seeding locker %s:\n%v", locker, err)
		}
	}

	logger.Info("Done")
}
