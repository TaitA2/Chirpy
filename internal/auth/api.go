package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	key := strings.TrimPrefix(authHeader, "ApiKey ")
	if len(key) < 1 {
		return "", fmt.Errorf("Invalid tokenString: %v", key)
	}
	return key, nil

}
