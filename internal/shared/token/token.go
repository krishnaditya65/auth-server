package token

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandom(
	length int,
) (string, error) {

	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(
		bytes,
	), nil
}
