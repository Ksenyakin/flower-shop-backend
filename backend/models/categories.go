package models

import (
	"flower-shop-backend/utils"
	"log"
	"time"
)

// Структура для категории
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Функция для создания новой категории
func CreateCategory(name string, description string) (int, error) {
	var categoryID int

	// SQL-запрос для вставки новой категории
	query := `
		INSERT INTO categories (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id;`

	// Выполнение запроса и получение id новой категории
	err := utils.DB.QueryRow(query, name, description).Scan(&categoryID)
	if err != nil {
		log.Println("Ошибка при создании категории:", err)
		return 0, err
	}

	return categoryID, nil
}
