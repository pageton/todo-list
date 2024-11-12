package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pageton/todo-list/config"
)

var jwtKey []byte

type TokenClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func InitConfig(cfg *config.Config) {
	if cfg == nil || cfg.SecretKey == "" {
		panic("JWT secret key cannot be empty")
	}
	jwtKey = []byte(cfg.SecretKey)
}

func CreateToken(userID string, username string) (string, error) {
	if userID == "" || username == "" {
		return "", errors.New("userID and username cannot be empty")
	}

	claims := &jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (*TokenClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token string cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}
