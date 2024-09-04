package main

import (
	"crypto/rand"
	"fmt"

	"github.com/zvdv/ECSS-Lockers/internal/crypto"
)

func main() {
	var keyByte [32]byte

	if _, err := rand.Read(keyByte[:]); err != nil {
		panic(err)
	}

	encoded := crypto.Base64.EncodeToString(keyByte[:])
	fmt.Println(encoded)
}
