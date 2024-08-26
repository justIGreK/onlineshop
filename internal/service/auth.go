package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"onlineshop/internal/models"
	grpcrequest "onlineshop/pkg/grpcReq"
)

const (
	tokenTTL   = 12 * time.Hour
	signingKey = "fsjklj235OIUJlknm24"
)

type Authorization interface {
	CreateUser(login, password string) (int, error)
	GetUser(login, password string) (models.User, error)
}
type tokenClaims struct {
	jwt.RegisteredClaims
	UserId   int    `json:"user_id"`
	UserRole string `json:"user_role"`
}

type AuthService struct {
	store Authorization
	grpcSender grpcrequest.GrpcRequest
}


func NewAuthService(store Authorization, grpcSender grpcrequest.GrpcRequest) *AuthService {
	return &AuthService{store: store, grpcSender: grpcSender}
}

func (s *AuthService) CreateUser(login string, password string, email string) (int, error) {
	isValid := s.grpcSender.GetRequest(email)
	if !isValid {
		return 0, errors.New("email is not valid")
	}
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
		}, user.Id, user.Role,
	})

	entryToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("error during generating token: %w", err)
	}
	return entryToken, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(accessToken *jwt.Token) (interface{}, error) {
		if _, ok := accessToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, "", fmt.Errorf("error during parsing token: %w", err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, claims.UserRole, nil
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash)
}
