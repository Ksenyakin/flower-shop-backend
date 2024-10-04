package handlers

import (
	"database/sql"
	"encoding/json"
	middlewares "flower-shop-backend/middleware"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"reflect"
	"time"
)

type JsonResponse struct {
	Message string `json:"message"`
}

// GetUserInfo получает информацию о пользователе
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Получение информации о пользователе")

	// Получаем токен из заголовка Authorization
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		logrus.Warn("Отсутствует заголовок Authorization")
		http.Error(w, `{"message": "Authorization header is required"}`, http.StatusUnauthorized)
		return
	}

	// Удаляем 'Bearer ' из токена, если он там есть
	if len(tokenStr) > len("Bearer ") {
		tokenStr = tokenStr[len("Bearer "):]
	} else {
		logrus.Warn("Неправильный формат токена")
		http.Error(w, `{"message": "Invalid token format"}`, http.StatusUnauthorized)
		return
	}

	// Парсим токен
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Warn("Неподдерживаемый метод подписи токена")
			return nil, nil
		}
		return []byte(getEnv("JWT_SECRET", "0000")), nil
	})

	if err != nil || !token.Valid {
		logrus.Warn("Неправильный токен: ", err)
		http.Error(w, `{"message": "Invalid token"}`, http.StatusUnauthorized)
		return
	}

	// Извлекаем данные из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		logrus.Warn("Неправильные данные в токене")
		http.Error(w, `{"message": "Invalid token"}`, http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		logrus.Warn("Ошибка получения email из токена")
		http.Error(w, `{"message": "Invalid token claims"}`, http.StatusUnauthorized)
		return
	}

	logrus.Info("Поиск пользователя с email: ", email)

	// Получаем пользователя по email
	user, err := models.GetUserByEmail(email)
	if err != nil {
		logrus.Error("Пользователь не найден: ", err)
		http.Error(w, `{"message": "User not found"}`, http.StatusNotFound)
		return
	}

	// Формируем ответ
	response := map[string]interface{}{
		"name":    user.Name,
		"phone":   user.Phone,
		"address": user.Address,
		"email":   user.Email,
	}

	logrus.Info("Информация о пользователе успешно получена")

	// Устанавливаем заголовок и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterUser регистрирует нового пользователя
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Регистрация нового пользователя")

	var user models.User

	middlewares.EnableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logrus.Warn("Ошибка декодирования данных запроса: ", err)
		http.Error(w, `{"message": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	// Проверка, существует ли уже пользователь с таким email
	var existingEmail string
	err := utils.DB.QueryRow(`SELECT email FROM users WHERE email = $1`, user.Email).Scan(&existingEmail)
	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		logrus.Warn("Пользователь с таким email уже существует: ", user.Email)
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		json.NewEncoder(w).Encode(JsonResponse{Message: "User already exists"})
		return
	} else if err != sql.ErrNoRows {
		logrus.Error("Ошибка проверки существования пользователя: ", err)
		http.Error(w, `{"message": "Failed to check user existence"}`, http.StatusInternalServerError)
		return
	}
	fmt.Println(reflect.TypeOf(user.Password))
	fmt.Println([]byte(user.Password))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	fmt.Println(hashedPassword)
	if err != nil {
		logrus.Error("Ошибка хэширования пароля: ", err)
		http.Error(w, `{"message": "Failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	// Вставка нового пользователя в базу данных
	_, err = utils.DB.Exec(`INSERT INTO users (email, password_hash, name, phone, address) VALUES ($1, $2, $3, $4, $5)`,
		user.Email, hashedPassword, user.Name, user.Phone, user.Address)
	if err != nil {
		logrus.Error("Ошибка регистрации пользователя: ", err)
		http.Error(w, `{"message": "Failed to register user"}`, http.StatusInternalServerError)
		return
	}

	logrus.Info("Пользователь успешно зарегистрирован: ", user.Email)

	// Успешная регистрация
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JsonResponse{Message: "Регистрация успешна"})
}

// LoginUser авторизует пользователя
func LoginUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Авторизация пользователя")

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
		logrus.Warn("Ошибка декодирования данных запроса: ", err)
		http.Error(w, `{"message": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	logrus.Info("Попытка входаfffff для email: ", loginData.Email)
	logrus.Info("Введенный пароль: ", loginData.Password)

	// Поиск пользователя по email и паролю
	user, err := models.GetUserByEmailAndPassword(loginData.Email, loginData.Password)
	if err != nil {
		logrus.Warn("Неудачная попытка входа для email: ", loginData.Email)
		http.Error(w, `{"message": "Invalid credentials"}`, http.StatusUnauthorized)
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		logrus.Error("Ошибка генерации токена: ", err)
		http.Error(w, `{"message": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	logrus.Info("Токен успешно сгенерирован для пользователя: ", user.Email)

	// Возвращение токена клиенту
	response := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Error("Ошибка отправки ответа: ", err)
		http.Error(w, `{"message": "Failed to send response"}`, http.StatusInternalServerError)
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
