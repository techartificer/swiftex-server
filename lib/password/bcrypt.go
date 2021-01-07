package password

import (
	"github.com/techartificer/swiftex/config"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generate hash password
func HashPassword(password string) (string, error) {
	bc := config.GetServer().BcryptCost
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bc)
	return string(bytes), err
}

// CheckPasswordHash compare raw and hash password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
