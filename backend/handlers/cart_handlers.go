package handlers

import (
	"encoding/json"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// AddToCart добавляет товар в корзину
func AddToCart(w http.ResponseWriter, r *http.Request) {
	var item models.CartItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Получение user_id из контекста запроса
	userID, ok := r.Context().Value("user_id").(float64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получение cart_id пользователя
	var cartID int
	err := utils.DB.QueryRow(`
        SELECT id FROM cart WHERE user_id = $1`, int(userID)).Scan(&cartID)
	if err != nil {
		http.Error(w, "Failed to get cart ID", http.StatusInternalServerError)
		return
	}

	// Добавление товара в корзину
	_, err = utils.DB.Exec(`
        INSERT INTO cart_items (cart_id, product_id, quantity)
        VALUES ($1, $2, $3)
        ON CONFLICT (cart_id, product_id) DO UPDATE
        SET quantity = cart_items.quantity + EXCLUDED.quantity`,
		cartID, item.ProductID, item.Quantity)
	if err != nil {
		http.Error(w, "Failed to add item to cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Item added to cart"})
}

// RemoveFromCart удаляет товар из корзины
func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID, err := strconv.Atoi(vars["cart_id"])
	if err != nil {
		http.Error(w, "Invalid cart ID", http.StatusBadRequest)
		return
	}
	productID, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	_, err = utils.DB.Exec(`
        DELETE FROM cart_items
        WHERE cart_id = $1 AND product_id = $2`,
		cartID, productID)
	if err != nil {
		http.Error(w, "Failed to remove item from cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Item removed from cart"})
}

// GetCart возвращает все товары в корзине пользователя
func GetCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	rows, err := utils.DB.Query(`
        SELECT ci.cart_id, ci.product_id, ci.quantity, p.name, p.price
        FROM cart_items ci
        JOIN products p ON ci.product_id = p.id
        WHERE ci.cart_id = (SELECT id FROM cart WHERE user_id = $1)`, userID)
	if err != nil {
		http.Error(w, "Failed to query cart", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.CartItem

	for rows.Next() {
		var item models.CartItem
		if err := rows.Scan(&item.CartID, &item.ProductID, &item.Quantity, &item.Name, &item.Price); err != nil {
			http.Error(w, "Failed to read cart items", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read cart items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
