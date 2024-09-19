package middlewares

import (
	"flower-shop-backend/utils"
	"net/http"
)

// AdminMiddleware проверяет, что пользователь является администратором
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Не удалось получить информацию о пользователе", http.StatusUnauthorized)
			return
		}

		isAdmin, err := utils.IsAdmin(userID)
		if err != nil || !isAdmin {
			http.Error(w, "Недостаточно прав для доступа", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
