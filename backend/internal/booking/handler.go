package booking

import (
	"fmt"
	"net/http"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/internal/seat"
	"cinema-booking/internal/ws"
	redisClient "cinema-booking/pkg/redis"

	"github.com/gin-gonic/gin"
)

type BookingRequest struct {
	UserID     string `json:"user_id"`
	SeatNumber string `json:"seat_number"`
	ShowtimeID string `json:"showtime_id"`
}

type ConfirmBookingRequest struct {
	ShowtimeID string `json:"showtime_id"`
}

func CreateBookingHandler(c *gin.Context) {

	var req BookingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	booking := models.Booking{
		UserID:     req.UserID,
		SeatNumber: req.SeatNumber,
		ShowtimeID: req.ShowtimeID,
		Status:     models.PENDING,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := CreateBooking(booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "booking created",
		"data":    booking,
	})
}

func ConfirmBookingHandler(hub *ws.Hub) gin.HandlerFunc {

	return func(c *gin.Context) {
		seatNumber := c.Param("seat_number")

		var req ConfirmBookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := ConfirmBooking(seatNumber, req.ShowtimeID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := seat.UpdateSeatStatus(req.ShowtimeID, seatNumber, "BOOKED"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		key := "seat_lock:" + seatNumber
		_ = redisClient.Client.Del(c.Request.Context(), key).Err()

		msg := fmt.Sprintf(`{
			"event":"seat_booked",
			"seat":"%s"
		}`, seatNumber)
		hub.Broadcast <- []byte(msg)

		c.JSON(http.StatusOK, gin.H{
			"message": "booking confirmed",
		})
	}

}
