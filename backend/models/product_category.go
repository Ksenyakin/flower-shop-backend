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

// AddProductCategory добавляет связь между товаром и категорией
func AddProductCategory(productID int, categoryID int) error {
	query := "INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2)"
	_, err := utils.DB.Exec(query, productID, categoryID)
	if err != nil {
		log.Println("Ошибка при добавлении категории к товару:", err)
		return err
	}
	return nil
}

// RemoveProductCategory удаляет связь между товаром и категорией
func RemoveProductCategory(productID int, categoryID int) error {
	query := "DELETE FROM product_categories WHERE product_id = $1 AND category_id = $2"
	_, err := utils.DB.Exec(query, productID, categoryID)
	if err != nil {
		log.Println("Ошибка при удалении категории товара:", err)
		return err
	}
	return nil
}

// GetCategoriesForProduct получает категории, связанные с товаром
func GetCategoriesForProduct(productID int) ([]Category, error) {
	query := `SELECT c.id, c.name FROM categories c 
              JOIN product_categories pc ON c.id = pc.category_id
              WHERE pc.product_id = $1`

	rows, err := utils.DB.Query(query, productID)
	if err != nil {
		log.Println("Ошибка при получении категорий товара:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			log.Println("Ошибка при чтении данных о категории:", err)
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
