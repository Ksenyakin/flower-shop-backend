package middlewares

import "net/http"

// EnableCORS добавляет заголовки для разрешения CORS-запросов
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем все источники
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Разрешаем определенные методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Разрешаем определенные заголовки (Content-Type, Authorization и другие)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Если запрос типа OPTIONS (предварительный запрос), отправляем пустой ответ
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Передаем запрос дальше по цепочке
		next.ServeHTTP(w, r)
	})
}
