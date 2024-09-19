package handlers

import (
	"database/sql"
	"encoding/json"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

// RegisterUser регистрирует нового пользователя
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Проверка, существует ли уже пользователь с таким email
	var existingEmail string
	err := utils.DB.QueryRow(`SELECT email FROM users WHERE email = $1`, user.Email).Scan(&existingEmail)
	if err == nil {
		// Если есть результат, значит пользователь с таким email уже существует
		http.Error(w, "User already exists", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		// Если произошла другая ошибка
		http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
		logrus.Error("Failed to check user existence: ", err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	logrus.Info("Hashed password: ", string(hashedPassword))

	// Вставка нового пользователя в базу данных
	_, err = utils.DB.Exec(`
        INSERT INTO users (email, password_hash, name, phone, address)
        VALUES ($1, $2, $3, $4, $5)`,
		user.Email, hashedPassword, user.Name, user.Phone, user.Address)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Успешная регистрация
	w.WriteHeader(http.StatusOK)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	type Config struct {
		Secret string
	}
	config := Config{
		Secret: getEnv("JWT_SECRET", "0000"),
	}
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmailAndPassword(loginData.Email, loginData.Password)
	if err != nil {
		logrus.Error("Login failed for email: ", loginData.Email, " Error: ", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := &utils.Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // Токен истекает через 1 час
			Issuer:    "your-app-name",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		logrus.Error("Failed to generate token: ", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
