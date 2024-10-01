package models

import (
	"database/sql"
	"errors"
	"flower-shop-backend/utils"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

// User представляет структуру пользователя
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	DayOfBithday string    `json:"birthday"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GetUserByEmailAndPassword проверяет учетные данные пользователя
func GetUserByEmailAndPassword(email, password string) (*User, error) {
	var user User

	// Получаем данные пользователя по email
	query := `SELECT id, email, password_hash, name, phone, address FROM users WHERE email = $1`
	row := utils.DB.QueryRow(query, email)

	// Сканируем данные в структуру User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Phone, &user.Address); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Проверяем, совпадает ли введенный пароль с хэшем из базы данных
	IsEqualPasswords(password, user.PasswordHash)

	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	// Подключение к базе данных (предполагается, что у вас есть глобальная переменная db)
	var user User
	query := "SELECT id, name, email, phone, address, birthday FROM users WHERE email = $1"
	row := utils.DB.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Address, &user.DayOfBithday)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		log.Println("Error querying user:", err)
		return nil, err // Ошибка при запросе
	}

	return &user, nil // Возвращаем найденного пользователя
}

func IsEqualPasswords(encryptedPassword string, expectedPassword string) error {
	bytesEncryptedPassword := []byte(encryptedPassword)
	bytesExpectedPassword := []byte(expectedPassword)
	err := bcrypt.CompareHashAndPassword(bytesEncryptedPassword, bytesExpectedPassword)
	return err
}
