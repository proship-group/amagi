package helpers

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	randSrc = rand.NewSource(time.Now().UnixNano())
)

const (
	rs6Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rs6LetterIdxBits = 6
	rs6LetterIdxMask = 1<<rs6LetterIdxBits - 1
	rs6LetterIdxMax  = 63 / rs6LetterIdxBits
)

// GenerateHashString generate a byte slice from string
func GenerateHashString(str string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(str), 10)
}

//RandString6 generate random strings from length
func RandString6(n int) string {
	b := make([]byte, n)
	cache, remain := randSrc.Int63(), rs6LetterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), rs6LetterIdxMax
		}
		idx := int(cache & rs6LetterIdxMask)
		if idx < len(rs6Letters) {
			b[i] = rs6Letters[idx]
			i--
		}
		cache >>= rs6LetterIdxBits
		remain--
	}
	return string(b)
}
