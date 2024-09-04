package internal

import (
	"github.com/zvdv/ECSS-Lockers/internal/env"
)

var (
	Domain string
)

func Initialize() {
	Domain = env.Env("DOMAIN")
}
