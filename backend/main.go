package main

import (
	"flower-shop-backend/routes" // Убедитесь, что путь правильный
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

func main() {
	r := routes.NewRouter()

	// Настройка CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Разрешить запросы с любого домена
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(r)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
