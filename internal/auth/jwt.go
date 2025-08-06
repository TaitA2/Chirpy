package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		Subject:   userID.String(),
	})
	tokenByte, err := (token.SignedString([]byte(tokenSecret)))
	if err != nil {
		return "", fmt.Errorf("Error signing token with secret string: %v", err)
	}
	return tokenByte, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error parsing token: %v", err)
	}
	id, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error getting userID: %v", err)
	}
	userID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error parsing uuid: %v", err)
	}
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if len(tokenString) < 1 {
		return "", fmt.Errorf("Invalid tokenString: %v", tokenString)
	}
	return tokenString, nil
}
