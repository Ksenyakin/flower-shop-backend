package models

import (
	"flower-shop-backend/utils"
	"log"
)

// Модель для связи товара с категорией
type ProductCategory struct {
	ProductID  int `json:"product_id"`
	CategoryID int `json:"category_id"`
}

// AddProductCategory добавляет категорию к товару, обновляя поле category_id в таблице products
func AddProductCategory(productID int, categoryID int) error {
	// Обновляем category_id у товара
	query := "UPDATE products SET category_id = $1 WHERE id = $2"
	_, err := utils.DB.Exec(query, categoryID, productID) // Исправлено: передаем categoryID, затем productID
	if err != nil {
		log.Println("Ошибка при добавлении категории к товару:", err)
		return err
	}
	return nil
}

// RemoveProductCategory удаляет категорию товара, устанавливая category_id в NULL
func RemoveProductCategory(productID int) error {
	// Устанавливаем category_id в NULL, чтобы "удалить" категорию товара
	query := "UPDATE products SET category_id = NULL WHERE id = $1"
	_, err := utils.DB.Exec(query, productID)
	if err != nil {
		log.Println("Ошибка при удалении категории товара:", err)
		return err
	}
	return nil
}

// GetCategoryForProduct получает категорию, связанную с товаром
func GetCategoryForProduct(productID int) (*Category, error) {
	// Получаем категорию, связанную с товаром через поле category_id
	query := `SELECT c.id, c.name, c.description FROM categories c
              JOIN products p ON c.id = p.category_id
              WHERE p.id = $1`

	var category Category
	err := utils.DB.QueryRow(query, productID).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		log.Println("Ошибка при получении категории товара:", err)
		return nil, err
	}

	return &category, nil
}
