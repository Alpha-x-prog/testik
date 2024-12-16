package handlers

import (
	"awesomeProject/db"
	"net/http"
	//	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateUser обновляет данные пользователя
func AddPet(c *gin.Context) {
	type Pet struct {
		OwnerID       int    `json:"owner_id" binding:"required"`
		Name          string `json:"name" binding:"required"`
		Age           int    `json:"age"`
		PetInfo       string `json:"pet_info"`
		SpecialNeeds  string `json:"special_needs"`
		Toy           int    `json:"toy"`
		Photos        int    `json:"photos"`
		ActivityID    int    `json:"activity_id"`
		ServicePerson string `json:"service_person"`
		Feed          int    `json:"feed"`
		Type          int    `json:"type"`
		Room          int    `json:"room"`
	}

	var newPet Pet
	if err := c.ShouldBindJSON(&newPet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	query := `INSERT INTO pets ( owner_id, name, age, pet_info, special_needs, toy, photos, activity_id, service_person, feed, type, room ) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12 ) RETURNING id`

	var petID int
	err := db.DB.QueryRow(query,
		newPet.OwnerID, newPet.Name, newPet.Age, newPet.PetInfo,
		newPet.SpecialNeeds, newPet.Toy, newPet.Photos, newPet.ActivityID,
		newPet.ServicePerson, newPet.Feed, newPet.Type, newPet.Room,
	).Scan(&petID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add pet", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pet added successfully", "pet_id": petID})
}
