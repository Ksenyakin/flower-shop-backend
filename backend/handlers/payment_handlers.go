package handlers

import (
	"encoding/json"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Обработчик процесса оплаты
func ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		logrus.Error("Некорректный запрос: ", err)
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	_, err := utils.DB.Exec(`
		INSERT INTO payments (order_id, payment_method, amount, status, transaction_id, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())`,
		payment.OrderID, payment.PaymentMethod, payment.Amount, payment.Status, payment.TransactionID)
	if err != nil {
		logrus.Error("Ошибка при обработке платежа: ", err)
		http.Error(w, "Ошибка при обработке платежа", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Платеж успешно обработан"})
}

// Обработчик обработки покупки пользователя
func ProcessPurchaseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		logrus.Error("Некорректный ID пользователя: ", err)
		http.Error(w, "Некорректный ID пользователя", http.StatusBadRequest)
		return
	}

	// Получаем сумму покупки из тела запроса
	var jsonData struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil {
		logrus.Error("Некорректный запрос: ", err)
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	// Обрабатываем покупку
	err = models.ProcessPurchase(userID, jsonData.Amount)
	if err != nil {
		logrus.Error("Ошибка при обработке покупки: ", err)
		http.Error(w, "Ошибка при обработке покупки", http.StatusInternalServerError)
		return
	}

	// Получаем обновленную информацию о пользователе
	user, err := models.GetUserByID(userID)
	if err != nil {
		logrus.Error("Ошибка при получении пользователя: ", err)
		http.Error(w, "Ошибка при получении пользователя", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"message":            "Покупка успешно обработана",
		"total_purchases":    user.TotalPurchases,
		"points":             user.Points,
		"loyalty_level":      user.LoyaltyLevel,
		"last_purchase_date": user.LastPurchaseDate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
