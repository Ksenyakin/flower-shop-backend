package handlers

import (
	"encoding/json"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"log"
	"net/http"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := utils.DB.Query("SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products")
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.ImageURL, &product.CreatedAt, &product.UpdatedAt); err != nil {
			http.Error(w, "Failed to read products", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	// Декодируем тело запроса в структуру Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println("Ошибка декодирования JSON:", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Вызываем функцию модели для добавления товара в базу данных
	if err := models.CreateProduct(&product); err != nil {
		log.Println("Ошибка добавления товара:", err)
		http.Error(w, "Не удалось добавить товар", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}
