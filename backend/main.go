package main

import (
	"flower-shop-backend/routes"
	"flower-shop-backend/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	// Инициализация базы данных
	utils.InitDB()

	// Создание нового роутера
	r := routes.SetupRoutes()

	logrus.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
