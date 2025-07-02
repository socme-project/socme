package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "default_state"
	}
	return base64.URLEncoding.EncodeToString(b)
}

func GenerateSecureToken(byteLength int) (string, error) {
	tokenBytes := make([]byte, byteLength)

	n, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to read random bytes for token: %w", err)
	}
	if n != byteLength {
		return "", fmt.Errorf("expected to read %d bytes for token, but read %d", byteLength, n)
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	return token, nil
}
