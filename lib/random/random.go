package random

import (
	"math/rand"
	"time"

	"github.com/techartificer/swiftex/lib/errors"
)

var isInitilaized = false

// Initialize set rand.Seed
func Initialize() {
	isInitilaized = true
	rand.Seed(time.Now().UnixNano())
}

var pool = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// GenerateRandomString a random string with len = l
func GenerateRandomString(size int) (string, error) {
	if !isInitilaized {
		return "", errors.NewError("Random not initialized")
	}
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[i] = pool[rand.Intn(len(pool))]
	}
	return string(bytes), nil
}
