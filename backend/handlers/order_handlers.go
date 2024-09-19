package handlers

import (
	"encoding/json"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// CartItem представляет элемент корзины
type CartItem struct {
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// CreateOrder обрабатывает создание нового заказа
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Получение user_id из контекста запроса
	userID, ok := r.Context().Value("user_id").(float64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Создание заказа
	result, err := utils.DB.Exec(`
        INSERT INTO orders (user_id, total_price, status)
        VALUES ($1, $2, 'pending') RETURNING id`,
		int(userID), order.TotalPrice)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get order ID", http.StatusInternalServerError)
		return
	}

	// Получение элементов корзины для заказа
	rows, err := utils.DB.Query(`
        SELECT p.id, p.name, ci.quantity, p.price
        FROM cart_items ci
        JOIN products p ON ci.product_id = p.id
        WHERE ci.cart_id = (
            SELECT id FROM cart WHERE user_id = $1
        )`, int(userID))
	if err != nil {
		http.Error(w, "Failed to query cart", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []CartItem

	for rows.Next() {
		var item CartItem
		if err := rows.Scan(&item.ProductID, &item.Name, &item.Quantity, &item.Price); err != nil {
			http.Error(w, "Failed to read cart items", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read cart items", http.StatusInternalServerError)
		return
	}

	// Добавление товаров в заказ
	for _, item := range items {
		_, err := utils.DB.Exec(`
            INSERT INTO order_items (order_id, product_id, quantity, price)
            VALUES ($1, $2, $3, $4)`,
			orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			http.Error(w, "Failed to add items to order", http.StatusInternalServerError)
			return
		}
	}

	// Очистка корзины пользователя
	_, err = utils.DB.Exec(`
        DELETE FROM cart_items
        WHERE cart_id = (
            SELECT id FROM cart WHERE user_id = $1
        )`, int(userID))
	if err != nil {
		http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"order_id": int(orderID)})
}

// GetOrder возвращает информацию о заказе по его ID
func GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["order_id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	rows, err := utils.DB.Query(`
        SELECT p.id, p.name, oi.quantity, p.price
        FROM order_items oi
        JOIN products p ON oi.product_id = p.id
        WHERE oi.order_id = $1`, orderID)
	if err != nil {
		http.Error(w, "Failed to query order", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []CartItem

	for rows.Next() {
		var item CartItem
		if err := rows.Scan(&item.ProductID, &item.Name, &item.Quantity, &item.Price); err != nil {
			http.Error(w, "Failed to read order items", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read order items", http.StatusInternalServerError)
		return
	}

	// Отправка данных в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
