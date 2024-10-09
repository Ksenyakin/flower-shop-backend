package models

import (
	"database/sql"
	"flower-shop-backend/utils"
	"log"
	"time"
)

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func CreateProduct(product *Product) error {
	query := "INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id"

	// Выполняем запрос и получаем ID созданного товара
	err := utils.DB.QueryRow(query, product.Name, product.Description, product.Price, product.Stock).Scan(&product.ID)
	if err != nil {
		log.Println("Ошибка при добавлении товара в базу данных:", err)
		return err
	}

	return nil
}

func DeleteProduct(productID int) error {
	query := "DELETE FROM products WHERE id = $1"

	// Выполняем запрос для удаления товара
	_, err := utils.DB.Exec(query, productID)
	if err != nil {
		log.Println("Ошибка при удалении товара:", err)
		return err
	}

	return nil
}

// GetProductByID находит товар по его ID
func GetProductByID(productID int) (*Product, error) {
	var product Product
	var imageURL sql.NullString // Для возможного NULL значения image_url

	query := "SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products WHERE id = $1"
	row := utils.DB.QueryRow(query, productID)

	// Сканируем строку результата запроса
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&imageURL, // Используем sql.NullString для поля, которое может быть NULL
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Товар не найден:", err)
			return nil, nil // Возвращаем nil, если товар не найден
		}
		log.Println("Ошибка при получении товара по ID:", err)
		return nil, err
	}

	// Присваиваем значение image_url, если оно не NULL
	if imageURL.Valid {
		product.ImageURL = imageURL.String
	} else {
		product.ImageURL = ""
	}

	return &product, nil
}
