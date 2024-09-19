package utils

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type Config struct {
	Secret string
}

// Claims структура для хранения данных в JWT
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// CreateToken создает новый JWT
func CreateToken(userID int, email string) (string, error) {
	config := Config{
		Secret: getEnv("JWT_SECRET", "0000"),
	}
	// Создание токена с данными пользователя и стандартными claims
	claims := Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Токен действует 24 часа
			Issuer:    "flowers-shop",
		},
	}

	// Создание нового токена с использованием алгоритма HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписывание токена
	signedToken, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseToken парсит и проверяет JWT
func ParseToken(tokenString string) (*Claims, error) {

	config := Config{
		Secret: getEnv("JWT_SECRET", "0000"),
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
