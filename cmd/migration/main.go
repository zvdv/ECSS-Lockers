package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/env"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Error.Fatal(err)
	}

	logger.Info.Println("DATABASE MIGRATION")

	database.Connect(fmt.Sprintf(
		"%s?authToken=%s",
		env.MustEnv("DATABASE_URL"),
		env.MustEnv("DATABASE_AUTH_TOKEN")))

	db, lock := database.Lock()
	defer lock.Unlock()

	schema, err := os.ReadFile("internal/database/schema.sql")
	if err != nil {
		logger.Error.Fatal(err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		logger.Error.Fatal(err)
	}

	logger.Info.Println("created schema.")

	logger.Info.Println("seeding 200 lockers..")
	// eeehhh i'm not proud of how this is being done but
	// database/sql does not support array type for query
	// arg out of the box :(
	for i := 0; i < 200; i++ {
		locker := fmt.Sprintf("ELW %03d", i+1)

		stmt, err := db.Prepare(`INSERT INTO locker (id) VALUES (:id);`)
		if err != nil {
			logger.Error.Fatal(err)
		}

		_, err = stmt.Exec(sql.Named("id", locker))
		if err != nil {
			logger.Error.Printf("error seeding locker %s:\n%v", locker, err)
		}
	}

	logger.Info.Println("Done")
}
