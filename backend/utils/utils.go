package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "default_state"
	}
	return base64.URLEncoding.EncodeToString(b)
}
