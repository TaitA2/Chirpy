package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAuth(t *testing.T) {
	pw := "changeme"
	hash, err := HashPassword(pw)
	if err != nil {
		t.Errorf("%v", err)
	}
	err = CheckPasswordhash(pw, hash)
	if err != nil {
		t.Errorf("Error checking pw: %v", err)
	}
}

func TestJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := time.Duration(5)
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("%v", err)
	} else {
		fmt.Printf("tokenString: %v\n", tokenString)
	}
	userID, err = ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Errorf("%v", err)
	}
}
