package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"onlineshop/internal/storage"
)

const (
	tokenTTL   = 12 * time.Hour
	signingKey = "fsjklj235OIUJlknm24"
)

type AuthService struct {
	store storage.Authorization
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func NewAuthService(store storage.Authorization) *AuthService {
	return &AuthService{store: store}
}

func (s *AuthService) CreateUser(login string, password string) (int, error) {
	password = generatePasswordHash(password)
	i, err := s.store.CreateUser(login, password)
	if err != nil {
		return 0, fmt.Errorf("error during creating acc: %w", err)
	}
	return i, nil
}

func (s *AuthService) GenerateToken(login, password string) (string, error) {
	user, err := s.store.GetUser(login, generatePasswordHash(password))
	if err != nil {
		return "", fmt.Errorf("error during getting user:%w", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}, user.Id,
	})

	entryToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("error during generating token: %w", err)
	}
	return entryToken, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(accessToken *jwt.Token) (interface{}, error) {
		if _, ok := accessToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("error during parsing token: %w", err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash)
}
