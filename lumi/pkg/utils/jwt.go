package utils

import (
	"errors"
	"time"

	"github.com/Mahaveer86619/lumi/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	UserID     uint   `json:"user_id"`
	Type       string `json:"type"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID uint) (string, string, error) {
	accessClaims := &JwtCustomClaims{
		UserID: userID,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 1 Day
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(config.GConfig.JWTSecret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &JwtCustomClaims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 Days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString([]byte(config.GConfig.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func ValidateToken(tokenString string, expectedType string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.GConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		if claims.Type != expectedType {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
