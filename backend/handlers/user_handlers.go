package handlers

import (
	"database/sql"
	"encoding/json"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

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
	user.PasswordHash = string(hashedPassword)

	_, err = utils.DB.Exec(`
        INSERT INTO users (email, password_hash, name, phone, address, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
		user.Email, user.PasswordHash, user.Name, user.Phone, user.Address)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	err := utils.DB.QueryRow("SELECT id, password_hash FROM users WHERE email = $1", credentials.Email).Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// For simplicity, we're not generating tokens here. Implement token generation if needed.

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
