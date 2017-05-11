package helpers

import "golang.org/x/crypto/bcrypt"

// GenerateHashString generate a byte slice from string
func GenerateHashString(str string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(str), 10)
}
