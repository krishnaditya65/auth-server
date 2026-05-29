package token

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandom(size int) (string, error) {
	b := make([]byte, size)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
