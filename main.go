package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	db        *sql.DB
	secretKey = "your_secret_key" // Замените на ваш секретный ключ
)

func main() {
	// Подключение к базе данных PostgreSQL
	var err error
	connStr := "user=postgres port=5433 password=qwerty dbname=productdb sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверка соединения
	if err = db.Ping(); err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Инициализация маршрутов
	router := gin.Default()

	router.POST("/login", loginHandler)
	router.GET("/welcome", authMiddleware, welcomeHandler)

	// Запуск сервера
	log.Println("Server running on http://localhost:8080")
	router.Run(":8080")
}

func loginHandler(c *gin.Context) {
	var creds struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	// Parse JSON input
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Query the database for the user
	var storedPassword string
	query := "SELECT password FROM user1 WHERE login = $1"
	err := db.QueryRow(query, creds.Login).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Compare passwords (you might want to use hashed passwords here)
	if creds.Password != storedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": creds.Login,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Return the token
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func authMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Next()
}

func welcomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome!"})
}
