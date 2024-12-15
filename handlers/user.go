package handlers

import (
	"net/http"
	//	"strconv"

	"awesomeProject/db"
	"github.com/gin-gonic/gin"
)

// GetUsers возвращает список всех пользователей
func GetUsers(c *gin.Context) {
	rows, err := db.DB.Query(`SELECT id, login FROM "user"`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var login string
		rows.Scan(&id, &login)
		users = append(users, map[string]interface{}{"id": id, "login": login})
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID возвращает пользователя по ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
	}
	err := db.DB.QueryRow(`SELECT id, login FROM "user" WHERE id = $1`, id).Scan(&user.ID, &user.Login)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser обновляет данные пользователя
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var creds struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := `UPDATE "user" SET login = $1, password = $2 WHERE id = $3`
	_, err := db.DB.Exec(query, creds.Login, creds.Password, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser удаляет пользователя
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM "user" WHERE id = $1`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}