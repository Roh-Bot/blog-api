package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Authentication struct {
	JWT        JWTAuthenticator
	Encryption Encryption
}

type JWTAuthenticator interface {
	GenerateToken(jwt.MapClaims) (string, error)
	ValidateToken(string) (bool, error)
}

type Encryption interface {
	Encrypt(data string) ([]byte, error)
}

func NewAuthentication(jwt2 *JWT, aes2 *AES) *Authentication {
	return &Authentication{
		JWT:        jwt2,
		Encryption: aes2,
	}
}
