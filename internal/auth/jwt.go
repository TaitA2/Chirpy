package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	tokenString, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", fmt.Errorf("Error signing token with secret string: %v", err)
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) { return tokenSecret, nil })
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error parsing token")
	}
	id, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error getting userID")
	}
	userID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error parsing uuid")
	}
	return userID, nil
}
