package utils

import "os"

// GetEnv возвращает значение переменной окружения или значение по умолчанию
func GetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

var JWTSecret = GetEnv("JWT_SECRET", "your_secret_key_here")
