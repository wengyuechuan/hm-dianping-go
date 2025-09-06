package utils

import (
	"errors"
	"hm-dianping-go/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims JWT声明
type Claims struct {
	UserID uint `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint) (string, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return "", errors.New("config not loaded")
	}

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWT.ExpireTime) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "hm-dianping",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return nil, errors.New("config not loaded")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}