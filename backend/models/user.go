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

// Константы для уровней лояльности
const (
	InitialLevel = "Начальный"
	SilverLevel  = "Серебряный"
	GoldLevel    = "Золотой"
)

// Пороги для перехода на следующий уровень
const (
	SilverThreshold = 5000.0  // Серебряный уровень при покупках от 5 000 рублей
	GoldThreshold   = 15000.0 // Золотой уровень при покупках от 15 000 рублей
)

// Процент начисления баллов в зависимости от уровня лояльности
var LevelPercentages = map[string]float64{
	InitialLevel: 0.03,
	SilverLevel:  0.06,
	GoldLevel:    0.10,
}

// Упорядоченный список уровней
var LevelOrder = []string{InitialLevel, SilverLevel, GoldLevel}

// User представляет структуру пользователя
type User struct {
	ID               int       `json:"id"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	PasswordHash     string    `json:"password_hash"`
	Name             string    `json:"name"`
	Phone            string    `json:"phone"`
	Address          string    `json:"address"`
	DayOfBirthday    string    `json:"birthday"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	TotalPurchases   float64   `json:"total_purchases"`
	Points           int       `json:"points"`
	LoyaltyLevel     string    `json:"loyalty_level"`
	LastPurchaseDate time.Time `json:"last_purchase_date"` // Новое поле
}

// Получение индекса уровня
func GetLevelIndex(level string) int {
	for i, l := range LevelOrder {
		if l == level {
			return i
		}
	}
	return 0 // По умолчанию начальный уровень
}

// Понижение уровня на один
func (u *User) DecreaseLevel(currentLevel string) string {
	idx := GetLevelIndex(currentLevel)
	if idx > 0 {
		idx--
	}
	return LevelOrder[idx]
}

// Расчет базового уровня на основе общей суммы покупок
func (u *User) CalculateBaseLevel() string {
	if u.TotalPurchases >= GoldThreshold {
		return GoldLevel
	} else if u.TotalPurchases >= SilverThreshold {
		return SilverLevel
	} else {
		return InitialLevel
	}
}

// Добавление новой покупки
func (u *User) AddPurchase(amount float64) error {
	// Обновляем общую сумму покупок
	u.TotalPurchases += amount

	// Обновляем дату последней покупки
	u.LastPurchaseDate = time.Now()

	// Обновляем уровень лояльности
	u.UpdateLoyaltyLevel()

	// Рассчитываем начисленные баллы
	percentage, exists := LevelPercentages[u.LoyaltyLevel]
	if !exists {
		return errors.New("неизвестный уровень лояльности")
	}

	pointsEarned := int(amount * percentage)
	u.Points += pointsEarned

	// Обновляем информацию в базе данных
	query := `
		UPDATE users SET
			total_purchases = $1,
			points = $2,
			loyalty_level = $3,
			last_purchase_date = $4,
			updated_at = NOW()
		WHERE id = $5
	`

	_, err := utils.DB.Exec(query, u.TotalPurchases, u.Points, u.LoyaltyLevel, u.LastPurchaseDate, u.ID)
	if err != nil {
		logrus.Error("Ошибка при обновлении данных пользователя: ", err)
		return err
	}

	logrus.Info("Покупка успешно обработана для пользователя: ", u.Email)
	return nil
}

// Обновление уровня лояльности с учетом неактивности
func (u *User) UpdateLoyaltyLevel() {
	// Определяем базовый уровень на основе общей суммы покупок
	baseLevel := u.CalculateBaseLevel()

	// Начинаем с базового уровня
	newLevel := baseLevel

	// Проверяем дату последней покупки
	if !u.LastPurchaseDate.IsZero() {
		inactiveDuration := time.Since(u.LastPurchaseDate)
		if inactiveDuration > 90*24*time.Hour { // Около 90 дней
			// Понижаем уровень на один
			newLevel = u.DecreaseLevel(baseLevel)
		}
	}

	u.LoyaltyLevel = newLevel
}

// Получение пользователя по email и паролю
func GetUserByEmailAndPassword(email, password string) (*User, error) {
	logrus.Info("Получение пользователя по email: ", email)

	var user User

	// Получаем данные пользователя по email
	query := `
		SELECT id, email, password_hash, name, phone, address, birthday,
		       total_purchases, points, loyalty_level, last_purchase_date,
		       created_at, updated_at
		FROM users WHERE email = $1
	`
	row := utils.DB.QueryRow(query, email)

	// Сканируем данные в структуру User
	if err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Phone,
		&user.Address, &user.DayOfBirthday, &user.TotalPurchases,
		&user.Points, &user.LoyaltyLevel, &user.LastPurchaseDate,
		&user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			logrus.Warn("Пользователь не найден: неверный email или пароль для email: ", email)
			return nil, errors.New("invalid email or password")
		}
		logrus.Error("Ошибка при получении данных пользователя: ", err)
		return nil, err
	}

	// Проверяем, совпадает ли введенный пароль с хэшем из базы данных
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		logrus.Warn("Пароли не совпадают для email: ", email)
		return nil, errors.New("invalid email or password")
	}

	logrus.Info("Пользователь успешно получен: ", user.Email)
	return &user, nil
}

// Получение пользователя по email
func GetUserByEmail(email string) (*User, error) {
	logrus.Info("Поиск пользователя по email: ", email)

	var user User
	query := `
		SELECT id, name, email, phone, address, birthday,
		       total_purchases, points, loyalty_level, last_purchase_date,
		       created_at, updated_at
		FROM users WHERE email = $1
	`
	row := utils.DB.QueryRow(query, email)

	err := row.Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Address,
		&user.DayOfBirthday, &user.TotalPurchases, &user.Points, &user.LoyaltyLevel,
		&user.LastPurchaseDate, &user.CreatedAt, &user.UpdatedAt,
	)

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

// Получение пользователя по ID
func GetUserByID(userID int) (*User, error) {
	logrus.Info("Поиск пользователя по ID: ", userID)

	var user User
	query := `
		SELECT id, name, email, phone, address, birthday,
		       total_purchases, points, loyalty_level, last_purchase_date,
		       created_at, updated_at
		FROM users WHERE id = $1
	`
	row := utils.DB.QueryRow(query, userID)

	err := row.Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Address,
		&user.DayOfBirthday, &user.TotalPurchases, &user.Points, &user.LoyaltyLevel,
		&user.LastPurchaseDate, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Warn("Пользователь не найден по ID: ", userID)
			return nil, nil // Пользователь не найден
		}
		logrus.Error("Ошибка при запросе пользователя: ", err)
		return nil, err // Ошибка при запросе
	}

	logrus.Info("Пользователь успешно найден: ", user.Email)
	return &user, nil // Возвращаем найденного пользователя
}

// Создание нового пользователя
func CreateUser(user *User) error {
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("Ошибка при хешировании пароля: ", err)
		return err
	}

	// Устанавливаем начальные значения
	user.PasswordHash = string(hashedPassword)
	user.LoyaltyLevel = InitialLevel
	user.TotalPurchases = 0
	user.Points = 0

	query := `
		INSERT INTO users (
			email, password_hash, name, phone, address, birthday,
			total_purchases, points, loyalty_level, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, NOW(), NOW()
		) RETURNING id
	`

	err = utils.DB.QueryRow(
		query,
		user.Email, user.PasswordHash, user.Name, user.Phone, user.Address,
		user.DayOfBirthday, user.TotalPurchases, user.Points, user.LoyaltyLevel,
	).Scan(&user.ID)

	if err != nil {
		logrus.Error("Ошибка при создании пользователя: ", err)
		return err
	}

	logrus.Info("Пользователь успешно создан: ", user.Email)
	return nil
}

// Проверка совпадения паролей
func IsEqualPasswords(hashedPassword []byte, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))
	logrus.Info("Проверка паролей.")

	if err != nil {
		logrus.Warn("Пароли не совпадают.")
	} else {
		logrus.Info("Пароли совпадают.")
	}
	return err
}

// Обработка покупки пользователя
func ProcessPurchase(userID int, purchaseAmount float64) error {
	// Получаем пользователя по ID
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	// Обрабатываем покупку
	err = user.AddPurchase(purchaseAmount)
	if err != nil {
		return err
	}

	logrus.Infof("Покупка на сумму %.2f обработана для пользователя %s", purchaseAmount, user.Email)
	return nil
}
