package middlewares

import (
	"flower-shop-backend/utils"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Secret string
}

// AuthMiddleware проверяет наличие и валидность JWT токена
func AuthMiddleware(next http.Handler) http.Handler {

	config := Config{
		Secret: getEnv("JWT_SECRET", "0000"),
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims := &utils.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
