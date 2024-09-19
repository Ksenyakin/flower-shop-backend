package handlers

import (
	"encoding/json"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := utils.DB.Exec(`
        INSERT INTO orders (user_id, total_price, status, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())`,
		order.UserID, order.TotalPrice, order.Status)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	rows, err := utils.DB.Query(`
        SELECT o.id, o.user_id, o.total_price, o.status, oi.product_id, p.name, oi.quantity, oi.price
        FROM orders o
        JOIN order_items oi ON o.id = oi.order_id
        JOIN products p ON oi.product_id = p.id
        WHERE o.id = $1`, orderID)
	if err != nil {
		http.Error(w, "Failed to fetch order", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var order models.Order
	var items []struct {
		ProductID int     `json:"product_id"`
		Name      string  `json:"name"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
	}

	for rows.Next() {
		var item struct {
			ProductID int
			Name      string
			Quantity  int
			Price     float64
		}
		if err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &item.ProductID, &item.Name, &item.Quantity, &item.Price); err != nil {
			http.Error(w, "Failed to read order items", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}
	order.Items = items

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
