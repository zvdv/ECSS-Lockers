package internal

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	hostemail    string
	hostpassword string
	mailserver   string
	mailport     int
	domain       string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("[WARN] .env file not found\n")
	}

	hostemail = envOrPanic("GMAIL_USER")
	hostpassword = envOrPanic("GMAIL_PASSWORD")
	mailserver = "smtp.gmail.com"
	mailport = 587
	domain = envOrPanic("ORIGIN")
}

func envOrPanic(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Printf("[WARN] env variable not set: %s \n", key)
	}
	return value
}
