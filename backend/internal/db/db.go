// backend/internal/db/db.go
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Config содержит параметры подключения к БД
type Config struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     string
	SSLMode  string
}

// InitDB инициализирует подключение к базе данных
func InitDB() *sql.DB {
	// Получаем параметры подключения из переменных окружения
	config := Config{
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "12345678"),
		DBName:   getEnv("DB_NAME", "FlowersShopBD"),
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Формируем строку подключения
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		config.User, config.Password, config.DBName, config.Host, config.Port, config.SSLMode)

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	log.Println("Успешное подключение к базе данных")
	return db
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
