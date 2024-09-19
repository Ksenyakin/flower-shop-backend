package utils

import (
	"fmt"
	_ "github.com/lib/pq"
)

// Проверяет, является ли пользователь администратором
func IsAdmin(userID int) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM admins WHERE user_id = $1", userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке прав администратора: %w", err)
	}
	return count > 0, nil
}
