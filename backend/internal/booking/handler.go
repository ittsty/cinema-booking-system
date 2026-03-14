package booking

import (
	"time"

	"cinema-booking/internal/models"

	"github.com/gin-gonic/gin"
)

type BookingRequest struct {
	UserID     string `json:"user_id"`
	SeatNumber string `json:"seat_number"`
	ShowtimeID string `json:"showtime_id"`
}

func CreateBookingHandler(c *gin.Context) {

	var req BookingRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	booking := models.Booking{
		UserID:     req.UserID,
		SeatNumber: req.SeatNumber,
		ShowtimeID: req.ShowtimeID,
		Status:     models.PENDING,
		CreatedAt:  time.Now(),
	}

	err := CreateBooking(booking)

	if err != nil {

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "booking created",
	})
}

func ConfirmBookingHandler(c *gin.Context) {

	seatNumber := c.Param("seat_number")

	err := ConfirmBooking(seatNumber)

	if err != nil {

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "booking confirmed",
	})
}
