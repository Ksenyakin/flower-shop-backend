package models

import (
	"flower-shop-backend/utils"
	"log"
)

type CartItemDetails struct {
	ID        int     `json:"id"`
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
}

func ViewCartDetails(userID int) ([]CartItemDetails, error) {
	query := `
		SELECT ci.id, p.id, p.name, p.price, ci.quantity, (p.price * ci.quantity) AS total
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = $1`

	rows, err := utils.DB.Query(query, userID)
	if err != nil {
		log.Println("Ошибка при получении корзины пользователя:", err)
		return nil, err
	}
	defer rows.Close()

	var cartItems []CartItemDetails
	for rows.Next() {
		var item CartItemDetails
		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.Name,
			&item.Price,
			&item.Quantity,
			&item.Total,
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
