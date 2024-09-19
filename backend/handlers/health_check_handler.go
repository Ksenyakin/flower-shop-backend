// handlers/health_check_handler.go
package handlers

import (
	"flower-shop-backend/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

// HealthCheckHandler проверяет соединение с базой данных
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Пробуем выполнить простой запрос к базе данных
	err := utils.DB.Ping()
	if err != nil {
		logrus.Errorf("Ошибка подключения к базе данных: %v", err)
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}

	// Если подключение успешно, возвращаем 200 OK
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Подключение к базе данных успешно"))
}
