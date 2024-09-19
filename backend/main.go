package main

import (
	middlewares "flower-shop-backend/middleware"
	"flower-shop-backend/routes"
	"flower-shop-backend/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	// Инициализация базы данных
	if err := utils.InitDB(); err != nil {
		logrus.Fatalf("Не удалось инициализировать базу данных: %v", err)
	}
	logrus.Info("База данных успешно инициализирована")

	r := routes.NewRouter()

	// Добавляем Middleware для авторизации
	r.Use(middlewares.AuthMiddleware)

	logrus.Info("Запуск сервера на порту :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logrus.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
