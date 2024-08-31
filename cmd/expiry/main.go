package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/email"
	"github.com/zvdv/ECSS-Lockers/internal/env"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Error.Fatal(err)
	}
	dbURL := fmt.Sprintf(
		"%s?authToken=%s",
		env.MustEnv("DATABASE_URL"),
		env.MustEnv("DATABASE_AUTH_TOKEN"))
	database.Connect(dbURL)
	email.Initialize()
}

// func queryExpiring() {

// }

func main() {
	//db, lock := database.Lock()
}
