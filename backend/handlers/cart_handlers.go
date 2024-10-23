package handlers

import (
	"encoding/json"
	"flower-shop-backend/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// Добавление товара в корзину
func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(r.URL.Query().Get("product_id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(r.URL.Query().Get("quantity"))
	if err != nil || quantity <= 0 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	err = models.AddToCart(userID, productID, quantity)
	if err != nil {
		http.Error(w, "Failed to add to cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Товар успешно добавлен в корзину"})
}

// Удаление товара из корзины
func RemoveFromCartHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем cart_item_id из параметров URL
	vars := mux.Vars(r)
	cartItemID, err := strconv.Atoi(vars["cart_item_id"])
	if err != nil {
		http.Error(w, "Invalid cart item ID", http.StatusBadRequest)
		return
	}

	// Удаление товара из корзины через модель
	err = models.RemoveFromCart(cartItemID)
	if err != nil {
		http.Error(w, "Failed to remove from cart", http.StatusInternalServerError)
		return
	}

	// Успешное удаление
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Товар успешно удален из корзины"})
}

// Обновление количества товара в корзине
func UpdateCartItemHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем cart_item_id из параметров URL
	vars := mux.Vars(r)
	cartItemID, err := strconv.Atoi(vars["cart_item_id"])
	if err != nil {
		http.Error(w, "Invalid cart item ID", http.StatusBadRequest)
		return
	}

	// Получаем количество товара из запроса
	quantity, err := strconv.Atoi(r.URL.Query().Get("quantity"))
	if err != nil || quantity <= 0 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	// Обновление количества товара в корзине через модель
	err = models.UpdateCartItem(cartItemID, quantity)
	if err != nil {
		http.Error(w, "Failed to update cart item", http.StatusInternalServerError)
		return
	}

	// Успешное обновление
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Количество товара успешно обновлено"})
}

// Просмотр корзины пользователя
func ViewCartHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из параметров URL
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Получаем корзину пользователя из модели
	cartItems, err := models.ViewCartDetails(userID)
	if err != nil {
		log.Println("Ошибка при получении корзины пользователя:", err)
		http.Error(w, "Failed to retrieve cart", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Возвращаем JSON с информацией о корзине
	json.NewEncoder(w).Encode(cartItems)
}
