package handlers

import (
	"encoding/json"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"net/http"
)

func ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := utils.DB.Exec(`
        INSERT INTO payments (order_id, payment_method, amount, status, transaction_id, created_at)
        VALUES ($1, $2, $3, $4, $5, NOW())`,
		payment.OrderID, payment.PaymentMethod, payment.Amount, payment.Status, payment.TransactionID)
	if err != nil {
		http.Error(w, "Failed to process payment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
