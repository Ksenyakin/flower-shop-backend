package models

import (
	"flower-shop-backend/utils"
	"log"
	"time"
)

type CartItem struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func AddToCart(userID, productID, quantity int) error {
	// SQL-запрос для добавления товара в корзину
	query := `
		INSERT INTO cart_items (user_id, product_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = NOW();`

	// Выполнение запроса
	_, err := utils.DB.Exec(query, userID, productID, quantity)
	if err != nil {
		log.Println("Ошибка при добавлении товара в корзину:", err)
		return err
	}

	log.Println("Товар успешно добавлен в корзину")
	return nil
}

func RemoveFromCart(cartItemID int) error {
	query := "DELETE FROM cart_items WHERE id = $1"
	_, err := utils.DB.Exec(query, cartItemID)
	if err != nil {
		log.Println("Ошибка при удалении товара из корзины:", err)
		return err
	}

	log.Println("Товар успешно удален из корзины")
	return nil
}

func UpdateCartItem(cartItemID, quantity int) error {
	query := `
		UPDATE cart_items
		SET quantity = $1, updated_at = NOW()
		WHERE id = $2`

	_, err := utils.DB.Exec(query, quantity, cartItemID)
	if err != nil {
		log.Println("Ошибка при обновлении количества товара в корзине:", err)
		return err
	}

	log.Println("Количество товара успешно обновлено")
	return nil
}

func ViewCart(userID int) ([]CartItem, error) {
	query := `
		SELECT ci.id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at
		FROM cart_items ci
		WHERE ci.user_id = $1`

	rows, err := utils.DB.Query(query, userID)
	if err != nil {
		log.Println("Ошибка при получении корзины пользователя:", err)
		return nil, err
	}
	defer rows.Close()

	var cartItems []CartItem
	for rows.Next() {
		var item CartItem
		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			log.Println("Ошибка при сканировании корзины:", err)
			return nil, err
		}
		cartItems = append(cartItems, item)
	}

	if err := rows.Err(); err != nil {
		log.Println("Ошибка при итерации по строкам корзины:", err)
		return nil, err
	}

	return cartItems, nil
}
