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
	rows, err := utils.DB.Query("SELECT id, category_id, name, description, price, stock, image_url, created_at, updated_at FROM products")
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
		var categoryID sql.NullInt32 // Для возможного NULL значения category_id
		var imageURL sql.NullString  // Для возможного NULL значения image_url

		if err := rows.Scan(
			&product.ID,
			&categoryID, // Используем sql.NullInt32 для category_id
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&imageURL, // Используем sql.NullString для image_url
			&product.CreatedAt,
			&product.UpdatedAt); err != nil {
			logrus.WithError(err).Error("Ошибка при чтении данных о продукте")
			http.Error(w, "Failed to read products", http.StatusInternalServerError)
			return
		}

		// Присваиваем значение category_id, если оно не NULL
		if categoryID.Valid {
			product.Category_id = int(categoryID.Int32) // Преобразуем к типу int
		} else {
			product.Category_id = 0 // Или любое другое значение по умолчанию
		}

		// Присваиваем значение image_url, если оно не NULL
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

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Получаем ID товара из параметров маршрута
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"]) // Конвертируем строку в число
	if err != nil {
		log.Println("Неверный ID товара:", err)
		http.Error(w, "Неверный ID товара", http.StatusBadRequest)
		return
	}

	var product models.Product

	// Декодируем тело запроса в структуру Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println("Ошибка декодирования JSON:", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Обновляем товар в базе данных
	if err := models.UpdateProduct(productID, &product); err != nil {
		log.Println("Ошибка обновления товара:", err)
		http.Error(w, "Не удалось обновить товар", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
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

// AddCategoryToProduct добавляет категорию к товару
func AddCategoryToProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["product_id"])
	categoryID, err2 := strconv.Atoi(vars["category_id"])

	if err != nil || err2 != nil {
		logrus.WithError(err).Error("Неверный ID товара или категории")
		http.Error(w, "Неверный ID товара или категории", http.StatusBadRequest)
		return
	}

	// Добавляем категорию к товару
	if err := models.AddProductCategory(productID, categoryID); err != nil {
		logrus.WithError(err).Error("Ошибка при добавлении категории к товару")
		http.Error(w, "Не удалось добавить категорию", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Категория добавлена к товару"})
}

// RemoveCategoryFromProduct удаляет категорию у товара
func RemoveCategoryFromProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["product_id"])

	if err != nil {
		logrus.WithError(err).Error("Неверный ID товара или категории")
		http.Error(w, "Неверный ID товара или категории", http.StatusBadRequest)
		return
	}

	// Удаляем категорию у товара
	if err := models.RemoveProductCategory(productID); err != nil {
		logrus.WithError(err).Error("Ошибка при удалении категории у товара")
		http.Error(w, "Не удалось удалить категорию", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Категория удалена у товара"})
}

// GetCategoriesForProduct получает категории для товара
func GetCategoriesForProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["product_id"])

	if err != nil {
		logrus.WithError(err).Error("Неверный ID товара")
		http.Error(w, "Неверный ID товара", http.StatusBadRequest)
		return
	}

	// Получаем категории для товара
	categories, err := models.GetCategoryForProduct(productID)
	if err != nil {
		logrus.WithError(err).Error("Ошибка при получении категорий для товара")
		http.Error(w, "Не удалось получить категории", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}
