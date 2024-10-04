package models

import (
	"database/sql"
	"errors"
	"flower-shop-backend/utils"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User представляет структуру пользователя
type User struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	PasswordHash  string    `json:"password_hash"`
	Name          string    `json:"name"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	DayOfBirthday string    `json:"birthday"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GetUserByEmailAndPassword проверяет учетные данные пользователя
func GetUserByEmailAndPassword(email, password string) (*User, error) {
	logrus.Info("Получение пользователя по email: ", email)

	var user User

	// Получаем данные пользователя по email
	query := `SELECT id, email, password_hash, name, phone, address FROM users WHERE email = $1`
	row := utils.DB.QueryRow(query, email)

	// Сканируем данные в структуру User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Phone, &user.Address); err != nil {
		if err == sql.ErrNoRows {
			logrus.Warn("Пользователь не найден: неверный email или пароль для email: ", email)
			return nil, errors.New("invalid email or password")
		}
		logrus.Error("Ошибка при получении данных пользователя: ", err)
		return nil, err
	}

	fmt.Print(user.PasswordHash) // Заебись, формат нужный

	// Проверяем, совпадает ли введенный пароль с хэшем из базы данных
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		logrus.Info("Проверка паролей: ", user.PasswordHash, " и ", password)
		a, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		fmt.Println(string(a))
		logrus.Warn("Пароли не совпадают для email: ", email)
		return nil, errors.New("invalid email or password")
	}

	logrus.Info("Пользователь успешно получен: ", user.Email)
	return &user, nil
}

// GetUserByEmail получает пользователя по email
func GetUserByEmail(email string) (*User, error) {
	logrus.Info("Поиск пользователя по email: ", email)

	var user User
	query := "SELECT id, name, email, phone, address, birthday FROM users WHERE email = $1"
	row := utils.DB.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Address, &user.DayOfBirthday)

	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Warn("Пользователь не найден по email: ", email)
			return nil, nil // Пользователь не найден
		}
		logrus.Error("Ошибка при запросе пользователя: ", err)
		return nil, err // Ошибка при запросе
	}

	logrus.Info("Пользователь успешно найден: ", user.Email)
	return &user, nil // Возвращаем найденного пользователя
}

// IsEqualPasswords проверяет совпадение паролей
func IsEqualPasswords(hashedPassword []byte, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))
	logrus.Info("Проверка паролей: ", hashedPassword, " и ", plainPassword)

	if err != nil {
		logrus.Warn("Пароли не совпадают.")
	} else {
		logrus.Info("Пароли совпадают.")
	}
	return err
}
