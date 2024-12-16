package handlers

import (
	"awesomeProject/db"
	"net/http"
	//	"strconv"

	"github.com/gin-gonic/gin"
)

type Room struct {
	Number          int    `json:"number"`
	StatusOccupancy bool   `json:"status_occupancy"`
	StatusCleaning  bool   `json:"status_cleaning"`
	Person          string `json:"person"`
	RoomType        int    `json:"room_type"`
	PricePerDay     int    `json:"price_per_day"`
	LiveStreamURL   string `json:"live_stream_url"`
	AvailableFrom   string `json:"available_from"` // В формате "YYYY-MM-DD"
	AvailableTo     string `json:"available_to"`   // В формате "YYYY-MM-DD"
}

func AddRoom(c *gin.Context) {

	var newRoom Room
	if err := c.ShouldBindJSON(&newRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := `INSERT INTO room ( number, status_occupancy, status_cleaning, person, room_type, price_per_day, live_stream_url, available_from, available_to ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := db.DB.Exec(query,
		newRoom.Number,
		newRoom.StatusOccupancy,
		newRoom.StatusCleaning,
		newRoom.Person,
		newRoom.RoomType,
		newRoom.PricePerDay,
		newRoom.LiveStreamURL,
		newRoom.AvailableFrom,
		newRoom.AvailableTo,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add room", "details": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ALL GOOD"})

}
