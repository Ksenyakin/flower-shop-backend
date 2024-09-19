package routes

import (
	"flower-shop-backend/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/api/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/api/login", handlers.LoginUser).Methods("POST")

	// Product routes
	r.HandleFunc("/api/products", handlers.GetProducts).Methods("GET")

	// Cart routes
	r.HandleFunc("/api/cart", handlers.GetCart).Methods("GET")
	r.HandleFunc("/api/cart/add", handlers.AddToCart).Methods("POST")
	r.HandleFunc("/api/cart/remove", handlers.RemoveFromCart).Methods("POST")

	// Order routes
	r.HandleFunc("/api/orders", handlers.CreateOrder).Methods("POST")
	r.HandleFunc("/api/orders/{id:[0-9]+}", handlers.GetOrder).Methods("GET")

	// Payment routes
	r.HandleFunc("/api/pay", handlers.ProcessPayment).Methods("POST")

	return r
}
