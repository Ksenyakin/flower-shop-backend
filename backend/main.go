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

	r := routes.NewRouter()

	// Добавляем Middleware для авторизации
	r.Use(AuthMiddleware)

	logrus.Info("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
