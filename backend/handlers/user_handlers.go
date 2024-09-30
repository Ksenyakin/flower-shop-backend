package handlers

import (
	"database/sql"
	"encoding/json"
	middlewares "flower-shop-backend/middleware"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// Получаем токен из заголовка Authorization
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	// Удаляем 'Bearer ' из токена, если он там есть
	if len(tokenStr) > len("Bearer ") {
		tokenStr = tokenStr[len("Bearer "):]
	} else {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	// Парсим токен
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil // Возвращаем nil, чтобы пройти проверку, что токен невалидный
		}
		return []byte(getEnv("JWT_SECRET", "0000")), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Извлекаем данные из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Получаем пользователя по email
	user, err := models.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Формируем ответ
	response := map[string]interface{}{
		"name":    user.Name,
		"phone":   user.Phone,
		"address": user.Address,
		"email":   user.Email,
	}

	// Устанавливаем заголовок и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterUser регистрирует нового пользователя
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	middlewares.EnableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Проверка, существует ли уже пользователь с таким email
	var existingEmail string
	err := utils.DB.QueryRow(`SELECT email FROM users WHERE email = $1`, user.Email).Scan(&existingEmail)
	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		// Если есть результат, значит пользователь с таким email уже существует
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		json.NewEncoder(w).Encode(map[string]string{"error": "User already exists"})
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Регистрация успешна"})
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

	// Декодирование входных данных
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Поиск пользователя по email и паролю
	user, err := models.GetUserByEmailAndPassword(loginData.Email, loginData.Password)
	if err != nil {
		logrus.Error("Login failed for email: ", loginData.Email, " Error: ", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Создание JWT токена
	claims := &utils.Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // Токен истекает через 1 час
			Issuer:    "your-app-name",
		},
	}

	// Подписывание токена с использованием секрета
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		logrus.Error("Failed to generate token: ", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Возвращение токена клиенту
	response := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Error("Failed to send response: ", err)
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
