package main

import (
	"database/sql"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/zvdv/ECSS-Lockers/internal/crypto"
	"github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/env"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/time"
)

type User struct {
	Name  string
	Email string
}

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Error.Fatal(err)
	}

	logger.Info.Println("DATABASE MIGRATION")

	database.Connect(fmt.Sprintf(
		"%s?authToken=%s",
		env.MustEnv("DATABASE_URL"),
		env.MustEnv("DATABASE_AUTH_TOKEN")))

	crypto.Initialize()

	db, lock := database.Lock()
	defer lock.Unlock()

	users := [...]User{
		{Name: "Matthew Johnson", Email: "matthewjohnson@uvic.ca"},
		{Name: "Lisa Thompson", Email: "lisathompson@uvic.ca"},
		{Name: "Caroline Singh", Email: "carolinesingh@uvic.ca"},
		{Name: "Alex Roberts", Email: "alexroberts@uvic.ca"},
		{Name: "Karan Patel", Email: "karanpatel@uvic.ca"},
		{Name: "David Anderson", Email: "davidanderson@uvic.ca"},
		{Name: "Henry Lee", Email: "henrylee@uvic.ca"},
		{Name: "Jessica Garcia", Email: "jessicagarcia@uvic.ca"},
		{Name: "Brian Clark", Email: "brianclark@uvic.ca"},
		{Name: "William Moore", Email: "williammoore@uvic.ca"},
		{Name: "Taylor Turner", Email: "taylorturner@uvic.ca"},
		{Name: "Zoe Martin", Email: "zoemartin@uvic.ca"},
		{Name: "Sam White", Email: "samwhite@uvic.ca"},
		{Name: "Yvonne Murphy", Email: "yvonnemurphy@uvic.ca"},
		{Name: "Austin Parker", Email: "austinparker@uvic.ca"},
		{Name: "Jonathan Hall", Email: "jonathanhall@uvic.ca"},
		{Name: "Michael Davis", Email: "michaeldavis@uvic.ca"},
		{Name: "Jason Sanchez", Email: "jasonsanchez@uvic.ca"},
		{Name: "Ryan Hill", Email: "ryanhill@uvic.ca"},
		{Name: "Paul Nelson", Email: "paulnelson@uvic.ca"},
	}

	for i, user := range users {
		locker := fmt.Sprintf("ELW %03d", i+1)
		exp := time.NextExpiryDate(time.Now())

		_, err := db.Exec(`
            INSERT INTO registration (locker, user, name, expiry)
            VALUES (:locker, :user, :name, :expiry);`,
			sql.Named("locker", locker),
			sql.Named("user", user.Email),
			sql.Named("name", user.Name),
			sql.Named("expiry", exp))

		if err != nil {
			panic(err)
		}
	}
}
