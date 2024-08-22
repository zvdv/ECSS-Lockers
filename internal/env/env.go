package env

import (
	"os"

	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

func MustEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		logger.Error.Fatalf("environment variable not set: %s\n", key)
	}
	return value
}

func Env(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		logger.Trace.Printf("environment variable not set: %s\n", key)
	}
	return value
}
