package handlers

import (
	"encoding/json"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

// AdminManageProducts управляет продуктами для администраторов
func AdminManageProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetProducts(w, r) // Можно переиспользовать существующий обработчик для получения продуктов
	case "POST":
		// Добавление нового продукта
		var product models.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Ошибка при декодировании запроса", http.StatusBadRequest)
			return
		}

		_, err := utils.DB.Exec("INSERT INTO products (name, description, price, stock, image_url) VALUES ($1, $2, $3, $4, $5)",
			product.Name, product.Description, product.Price, product.Stock, product.ImageURL)
		if err != nil {
			http.Error(w, "Не удалось добавить продукт", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		logrus.Info("Продукт успешно добавлен")
	case "PUT":
		// Обновление существующего продукта
		var product models.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Ошибка при декодировании запроса", http.StatusBadRequest)
			return
		}

		_, err := utils.DB.Exec("UPDATE products SET name = $1, description = $2, price = $3, stock = $4, image_url = $5 WHERE id = $6",
			product.Name, product.Description, product.Price, product.Stock, product.ImageURL, product.ID)
		if err != nil {
			http.Error(w, "Не удалось обновить продукт", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		logrus.Info("Продукт успешно обновлен")
	case "DELETE":
		vars := mux.Vars(r)
		productID := vars["id"]

		_, err := utils.DB.Exec("DELETE FROM products WHERE id = $1", productID)
		if err != nil {
			http.Error(w, "Не удалось удалить продукт", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		logrus.Info("Продукт успешно удален")
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
