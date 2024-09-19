package routes

import (
	"flower-shop-backend/handlers"
	"flower-shop-backend/middleware"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Пользователи
	r.HandleFunc("/api/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/api/login", handlers.LoginUser).Methods("POST")

	// Товары
	r.HandleFunc("/api/products", handlers.GetProducts).Methods("GET")

	// Корзина
	r.HandleFunc("/api/cart/{user_id:[0-9]+}", handlers.GetCart).Methods("GET")
	r.HandleFunc("/api/cart/add", handlers.AddToCart).Methods("POST")
	r.HandleFunc("/api/cart/remove/{cart_id:[0-9]+}/{product_id:[0-9]+}", handlers.RemoveFromCart).Methods("DELETE")

	// Заказы
	r.HandleFunc("/api/orders", handlers.CreateOrder).Methods("POST")
	r.HandleFunc("/api/orders/{order_id:[0-9]+}", handlers.GetOrder).Methods("GET")

	// Платежи
	r.HandleFunc("/api/pay", handlers.ProcessPayment).Methods("POST")

	// Добавляем Middleware для авторизации
	r.Use(middlewares.AuthMiddleware)

	return r
}
