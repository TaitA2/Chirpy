package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("Error getting random bytes for refresh token: %v", err)
	}
	src := []byte(key)
	encoded := hex.EncodeToString(src)
	return encoded, nil
}
