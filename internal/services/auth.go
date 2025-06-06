package services

import (
	"github.com/Roh-Bot/blog-api/internal/auth"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/Roh-Bot/blog-api/internal/store"
	"github.com/Roh-Bot/blog-api/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"slices"

	"time"
)

type AuthService struct {
	config *config.AtomicConfig
	logger logger.Logger
	store  store.Store
	auth   *auth.Authentication
}

func (a *AuthService) GenerateToken(userId string) (token string, err error) {
	encryptedUserId, err := a.auth.Encryption.Encrypt(userId)
	if err != nil {
		return
	}

	//Generate token for user
	tokenTTL := time.Duration(a.config.Get().Auth.TokenTTL) * time.Minute

	tokenClaims := jwt.MapClaims{
		"iss": a.config.Get().Auth.Issuer,
		"sub": string(encryptedUserId),
		"aud": a.config.Get().Auth.Audience,
		"exp": time.Now().Add(tokenTTL).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}
	token, err = a.auth.JWT.GenerateToken(tokenClaims)
	if err != nil {
		return
	}
	return
}

// IsValid is used to validate the user
func (a *AuthService) IsValid(username string) bool {
	return slices.Contains(a.config.Get().Auth.ValidUsers, username)
}

func (a *AuthService) ValidateToken(token string) (bool, error) {
	return a.auth.JWT.ValidateToken(token)
}
