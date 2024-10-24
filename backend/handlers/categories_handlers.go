package handlers

import (
	"encoding/json"
	"flower-shop-backend/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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
func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID категории из параметров маршрута
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Удаляем категорию
	if err := models.DeleteCategory(categoryID); err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Категория успешно удалена"))
}
func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID категории из параметров маршрута
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var req models.Category

	// Декодируем тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Обновляем категорию
	if err := models.UpdateCategory(categoryID, req.Name, req.Description); err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Категория успешно обновлена"))
}
