package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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

func Hash(
	token string,
) string {

	sum := sha256.Sum256(
		[]byte(token),
	)

	return hex.EncodeToString(
		sum[:],
	)
}
