package lib

import (
	"crypto/rand"
	"encoding/hex"
)

// ? Function that generate a random string (UUID, State)
func GenerateRandomID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
