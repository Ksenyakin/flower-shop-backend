package routes

import (
	"flower-shop-backend/handlers"
	"flower-shop-backend/middleware"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Пользователи
	r.HandleFunc("/api/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/api/login", handlers.LoginUser).Methods("POST")
	r.HandleFunc("/api/userinfo", handlers.GetUserInfo).Methods("GET")

	// Товары
	r.HandleFunc("/api/products", handlers.GetProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", handlers.GetProductByID).Methods("GET")
	r.HandleFunc("/api/addProduct", handlers.AddProduct).Methods("POST")
	r.HandleFunc("/api/products/{id}", handlers.DeleteProduct).Methods("DELETE")

	// Категории Товаров
	r.HandleFunc("/api/products/{product_id}/categories", handlers.GetCategoriesForProduct).Methods("GET")
	r.HandleFunc("/api/products/{product_id}/categories/{category_id}", handlers.AddCategoryToProduct).Methods("POST")
	r.HandleFunc("/api/products/{product_id}/categories/{category_id}", handlers.RemoveCategoryFromProduct).Methods("DELETE")

	// Корзина
	r.HandleFunc("/api/cart/{user_id:[0-9]+}", handlers.GetCart).Methods("GET")
	r.HandleFunc("/api/cart/add", handlers.AddToCart).Methods("POST")
	r.HandleFunc("/api/cart/remove/{cart_id:[0-9]+}/{product_id:[0-9]+}", handlers.RemoveFromCart).Methods("DELETE")

	// Заказы
	r.HandleFunc("/api/orders", handlers.CreateOrder).Methods("POST")
	r.HandleFunc("/api/orders/{order_id:[0-9]+}", handlers.GetOrder).Methods("GET")

	// Платежи
	r.HandleFunc("/api/pay", handlers.ProcessPayment).Methods("POST")
	r.HandleFunc("/api/purchase/{id}", handlers.ProcessPurchaseHandler).Methods("POST")

	adminRoutes := r.PathPrefix("/api/admin").Subrouter()
	adminRoutes.Use(middlewares.AdminMiddleware)
	adminRoutes.HandleFunc("/products", handlers.AdminManageProducts).Methods("GET", "POST", "PUT", "DELETE")

	return r
}
