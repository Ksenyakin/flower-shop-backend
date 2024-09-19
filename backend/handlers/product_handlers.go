package handlers

import (
	"encoding/json"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
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
