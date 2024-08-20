package env

import (
	"fmt"
	"os"

	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

func MustEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		logger.Fatal(fmt.Sprintf("environment variable not set: %s", key))
	}
	return value
}
