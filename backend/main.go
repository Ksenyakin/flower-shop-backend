package main

import (
	"flower-shop-backend/internal/db"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Инициализация базы данных
	database := db.InitDB()
	defer database.Close()

	// Остальной код сервера...
}
