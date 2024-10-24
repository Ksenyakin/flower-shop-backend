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

// DeleteCategory удаляет категорию по её ID
func DeleteCategory(categoryID int) error {
	// SQL-запрос для удаления категории
	query := "DELETE FROM categories WHERE id = $1"

	// Выполнение запроса
	_, err := utils.DB.Exec(query, categoryID)
	if err != nil {
		log.Println("Ошибка при удалении категории:", err)
		return err
	}

	return nil
}

// UpdateCategory обновляет данные категории по её ID
func UpdateCategory(categoryID int, name string, description string) error {
	// SQL-запрос для обновления категории
	query := `
		UPDATE categories
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3`

	// Выполнение запроса
	_, err := utils.DB.Exec(query, name, description, categoryID)
	if err != nil {
		log.Println("Ошибка при обновлении категории:", err)
		return err
	}

	return nil
}
