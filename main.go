package main

import (
	"log"

	"awesomeProject/db"
	"awesomeProject/handlers"
	"awesomeProject/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	// Подключение к базе данных
	db.ConnectDatabase()
	defer db.CloseDatabase()

	// Инициализация маршрутов
	router := gin.Default()

	router.POST("/register", handlers.RegisterUser)
	router.POST("/login", handlers.LoginHandler)
	router.GET("/users", middlewares.AuthMiddleware, handlers.GetUsers)
	router.GET("/users/:id", middlewares.AuthMiddleware, handlers.GetUserByID)
	router.PUT("/users/:id", middlewares.AuthMiddleware, handlers.UpdateUser)
	router.DELETE("/users/:id", middlewares.AuthMiddleware, handlers.DeleteUser)

	log.Println("Server running on http://localhost:8080")
	router.Run(":8080")
}
