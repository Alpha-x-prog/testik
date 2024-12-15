package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// ConnectDatabase подключается к базе данных
func ConnectDatabase() {
	var err error
	connStr := "user=postgres port=5433 password=qwerty dbname=productdb sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Проверка соединения
	if err = DB.Ping(); err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	log.Println("Database connected successfully.")
}

// CloseDatabase закрывает соединение с базой данных
func CloseDatabase() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed.")
	}
}
