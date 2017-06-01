package auth

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandBytesToken generate random bytes token
func GenerateRandBytesToken(length int) []byte {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil
	}

	return b
}

// GenerateRandStringToken generate random string base token
func GenerateRandStringToken(length int) string {
	return base64.URLEncoding.EncodeToString(GenerateRandBytesToken(length))
}
