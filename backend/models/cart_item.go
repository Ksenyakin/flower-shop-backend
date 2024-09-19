package models

// CartItem представляет элемент корзины
type CartItem struct {
	CartID    int     `json:"cart_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Name      string  `json:"name,omitempty"`
	Price     float64 `json:"price,omitempty"`
}
