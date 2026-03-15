package booking

import (
	"fmt"
	"net/http"
	"time"

	"cinema-booking/internal/audit"
	"cinema-booking/internal/models"
	"cinema-booking/internal/mq"
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
	UserID     string `json:"user_id"`
}

func CreateBookingHandler(c *gin.Context) {

	var req BookingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = audit.LogEvent("SYSTEM_ERROR", req.UserID, req.SeatNumber, req.ShowtimeID, err.Error())

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
		_ = audit.LogEvent("SYSTEM_ERROR", req.UserID, req.SeatNumber, req.ShowtimeID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	_ = seat.UpdateSeatStatus(req.ShowtimeID, req.SeatNumber, "LOCKED")
	_ = audit.LogEvent("BOOKING_CREATED", req.UserID, req.SeatNumber, req.ShowtimeID, "booking created successfully")

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
			_ = audit.LogEvent("SYSTEM_ERROR", "", seatNumber, req.ShowtimeID, err.Error())

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := ConfirmBooking(seatNumber, req.ShowtimeID); err != nil {
			_ = audit.LogEvent("SYSTEM_ERROR", "", seatNumber, req.ShowtimeID, err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := seat.UpdateSeatStatus(req.ShowtimeID, seatNumber, "BOOKED"); err != nil {
			_ = audit.LogEvent("SYSTEM_ERROR", "", seatNumber, req.ShowtimeID, err.Error())

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

		_ = audit.LogEvent("BOOKING_SUCCESS", "", seatNumber, req.ShowtimeID, "booking confirmed successfully")

		event := mq.BookingSuccessEvent{
			UserID:     req.UserID,
			SeatNumber: seatNumber,
			ShowtimeID: req.ShowtimeID,
			Status:     "BOOKED",
			CreatedAt:  time.Now(),
		}

		if err := mq.PublishBookingSuccess(event); err != nil {
			_ = audit.LogEvent("SYSTEM_ERROR", req.UserID, seatNumber, req.ShowtimeID, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "booking confirmed",
		})
	}

}
