package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), 1)
	if err != nil {
		return "", fmt.Errorf("Error hashing password: %v", err)
	}
	return string(hash), nil
}

func CheckPasswordhash(pw, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}
