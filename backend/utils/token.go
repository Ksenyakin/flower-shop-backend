package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

// Secret key for JWT
var jwtSecret = []byte("your_secret_key")

// ParseToken проверяет JWT токен и возвращает userID
func ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Подтверждаем, что алгоритм токена верный
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id claim")
	}

	return int(userID), nil
}
