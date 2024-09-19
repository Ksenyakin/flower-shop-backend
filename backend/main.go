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

	// Настройка маршрутов

	// Добавляем Middleware для авторизации
	// Middleware может быть применен глобально или на уровне определённых маршрутов
	//r.Use(middlewares.AuthMiddleware)
	// Для демонстрации используем глобальное применение (например, только для защищённых маршрутов)
	//r.PathPrefix("/protected").Subrouter().Use(middlewares.AuthMiddleware)

	// Настройка CORS (если необходимо)
	// corsHandler := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	// 	AllowedHeaders: []string{"Content-Type", "Authorization"},
	// }).Handler
	// http.Handle("/", corsHandler(r))

	logrus.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
