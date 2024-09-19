package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB инициализирует соединение с базой данных
func InitDB() error {
	var err error
	connStr := "user=postgres dbname=FlowersShopBD sslmode=disable password=12345678"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}
	return nil
}
