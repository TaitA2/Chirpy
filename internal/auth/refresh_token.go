package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
)

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	data, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("Error getting random bytes for refresh token: %v", err)
	}
	src := []byte(strconv.Itoa(data))
	encoded := hex.EncodeToString(src)
	return encoded, nil
}
