package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	ErrorsIllegalSign  = errors.New("非法签名算法")
	ErrorsExpiredToken = errors.New("token 已过期")
	ErrorsInvalidToken = errors.New("token 无效")
)

var jwtKey = []byte("bT7@kL2#xV9!mQ4$rN8zC1&dF6pY3wHsJ5uE0tR2yI8oP4aS7dG9hK1lZ3cX6vBn")

// GetToken 生成token
func GetToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	}

	//生成一个未签名的token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//签名并返回
	return token.SignedString(jwtKey)
}

// ParseToken 校验token
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrorsIllegalSign
		}
		return jwtKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrorsExpiredToken
		}
		return nil, err
	}

	if !token.Valid {
		return nil, ErrorsInvalidToken
	}

	return claims, nil
}
