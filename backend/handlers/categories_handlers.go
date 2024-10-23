package handlers

import (
	"encoding/json"
	"flower-shop-backend/models"
	"log"
	"net/http"
)

// Структура для запроса создания категории
type CreateCategoryRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Хендлер для создания новой категории
func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest

	// Декодирование JSON-запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Вызов функции для создания категории
	categoryID, err := models.CreateCategory(req.Name, req.Description)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		log.Println("Ошибка при создании категории:", err)
		return
	}

	// Ответ с id новой категории
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Категория успешно создана",
		"id":      categoryID,
	})
}
