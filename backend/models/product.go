package models

import (
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
