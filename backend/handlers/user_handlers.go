package handlers

import (
	"encoding/json"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// RegisterUser регистрирует нового пользователя
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = utils.DB.Exec(`
        INSERT INTO users (email, password_hash, name, phone, address)
        VALUES ($1, $2, $3, $4, $5)`,
		user.Email, hashedPassword, user.Name, user.Phone, user.Address)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// LoginUser аутентифицирует пользователя и возвращает JWT
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var storedHash string
	err := utils.DB.QueryRow(`
        SELECT password_hash FROM users WHERE email = $1`,
		user.Email).Scan(&storedHash)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.PasswordHash)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
