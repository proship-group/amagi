package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"io/ioutil"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
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

// GenerateHMACToken generate hmac token from jwt-go
func GenerateHMACToken(token *jwt.Token, signKey *rsa.PrivateKey) (string, error) {
	return token.SignedString(signKey)
}

// InitPrivateKeys initialize private/public keys for jwt
func InitPrivateKeys(signKey *rsa.PrivateKey, verifyKey *rsa.PublicKey) {
	signBytes, err := ioutil.ReadFile(os.Getenv("PRIV_KEY_PATH"))
	if err != nil {
		panic(err)
	}
	signkey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}
	(*signKey) = *signkey

	verifyBytes, err := ioutil.ReadFile(os.Getenv("PUB_KEY_PATH"))
	if err != nil {
		panic(err)
	}
	verifykey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}

	(*verifyKey) = *verifykey
}
