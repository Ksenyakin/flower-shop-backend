package handlers

import (
	"database/sql"
	"encoding/json"
	"flower-shop-backend/models"
	"flower-shop-backend/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
)

// GetProducts возвращает список продуктов
func GetProducts(w http.ResponseWriter, r *http.Request) {
	// Выполняем запрос к базе данных
	rows, err := utils.DB.Query("SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products")
	if err != nil {
		logrus.WithError(err).Error("Ошибка выполнения запроса к базе данных")
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product

	// Проходим по строкам результата запроса
	for rows.Next() {
		var product models.Product
		var imageURL sql.NullString // Для возможного NULL значения image_url

		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&imageURL, // Используем sql.NullString для поля, которое может быть NULL
			&product.CreatedAt,
			&product.UpdatedAt); err != nil {
			logrus.WithError(err).Error("Ошибка при чтении данных о продукте")
			http.Error(w, "Failed to read products", http.StatusInternalServerError)
			return
		}

		// Присваиваем значение из sql.NullString к полю структуры
		if imageURL.Valid {
			product.ImageURL = imageURL.String
		} else {
			product.ImageURL = "" // Или другое значение по умолчанию, если URL отсутствует
		}

		products = append(products, product)
	}

	// Проверка на ошибки, возникшие во время итерации
	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error("Ошибка при итерации по строкам результата запроса")
		http.Error(w, "Failed to read products", http.StatusInternalServerError)
		return
	}

	// Возвращаем список продуктов
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		logrus.WithError(err).Error("Ошибка при кодировании JSON ответа")
		http.Error(w, "Failed to encode products", http.StatusInternalServerError)
	}
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	// Декодируем тело запроса в структуру Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println("Ошибка декодирования JSON:", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Вызываем функцию модели для добавления товара в базу данных
	if err := models.CreateProduct(&product); err != nil {
		log.Println("Ошибка добавления товара:", err)
		http.Error(w, "Не удалось добавить товар", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Получаем ID товара из параметров маршрута
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"]) // Конвертируем строку в число
	if err != nil {
		log.Println("Неверный ID товара:", err)
		http.Error(w, "Неверный ID товара", http.StatusBadRequest)
		return
	}

	// Вызываем функцию модели для удаления товара
	if err := models.DeleteProduct(productID); err != nil {
		log.Println("Ошибка при удалении товара:", err)
		http.Error(w, "Не удалось удалить товар", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Товар успешно удалён"))
}

// GetProductByID обрабатывает запрос для получения товара по его ID
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	// Получаем ID товара из URL
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"]) // Конвертируем строку в число
	if err != nil {
		logrus.WithError(err).Error("Неверный ID товара")
		http.Error(w, "Неверный ID товара", http.StatusBadRequest)
		return
	}

	// Вызываем функцию модели для получения товара по ID
	product, err := models.GetProductByID(productID)
	if err != nil {
		logrus.WithError(err).Error("Ошибка при получении товара по ID")
		http.Error(w, "Не удалось получить товар", http.StatusInternalServerError)
		return
	}

	// Если товар не найден
	if product == nil {
		http.Error(w, "Товар не найден", http.StatusNotFound)
		return
	}

	// Возвращаем товар в формате JSON
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		logrus.WithError(err).Error("Ошибка при кодировании JSON ответа")
		http.Error(w, "Не удалось закодировать ответ", http.StatusInternalServerError)
	}
}
