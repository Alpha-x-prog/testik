package handlers

import (
	"net/http"
	"time"

	"awesomeProject/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = "your_secret_key" // Замените на ваш секретный ключ

// RegisterUser регистрирует нового пользователя
func RegisterUser(c *gin.Context) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		Phone    int    `json:"phone"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	queryUserList := `INSERT INTO userlist (email, password) VALUES ($1, $2) RETURNING id`
	var userID int
	err := db.DB.QueryRow(queryUserList, creds.Email, creds.Password).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error 1"})
		return
	}

	queryPetOwner := `INSERT INTO "petOwner" (id, name, surname, phone, email, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.DB.Exec(queryPetOwner, userID, creds.Name, creds.Surname, creds.Phone, creds.Email, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// LoginHandler обрабатывает логин пользователя
func LoginHandler(c *gin.Context) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var storedPassword string
	query := `SELECT password FROM userList WHERE email = $1`
	err := db.DB.QueryRow(query, creds.Email).Scan(&storedPassword)
	if err != nil || creds.Password != storedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": creds.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
