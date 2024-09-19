package models

import (
	"database/sql"
	"errors"
	"flower-shop-backend/utils"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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

	// Логируем хеш пароля и введенный пароль
	logrus.Infof("Password hash from database: %s", user.PasswordHash)
	logrus.Infof("Entered password: %s", password)

	// Проверяем, совпадает ли введенный пароль с хэшем из базы данных
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		logrus.Warning("Password mismatch for email: ", email)
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}
