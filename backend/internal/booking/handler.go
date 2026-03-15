package booking

import (
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

	if req.UserID == "" || req.SeatNumber == "" || req.ShowtimeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id, seat_number and showtime_id are required",
		})
		return
	}

	seatData, err := seat.GetSeatByShowtimeAndNumber(req.ShowtimeID, req.SeatNumber)
	if err != nil {
		_ = audit.LogEvent("SYSTEM_ERROR", req.UserID, req.SeatNumber, req.ShowtimeID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if seatData.Status != "LOCKED" {
		c.JSON(http.StatusConflict, gin.H{
			"message": "seat is not in LOCKED state",
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
	exists, err := HasPendingBooking(req.SeatNumber, req.ShowtimeID)
	if err != nil {
		_ = audit.LogEvent("SYSTEM_ERROR", req.UserID, req.SeatNumber, req.ShowtimeID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"message": "pending booking already exists for this seat",
		})
		return
	}
	if err := CreateBooking(booking); err != nil {
		_ = audit.LogEvent("SYSTEM_ERROR", req.UserID, req.SeatNumber, req.ShowtimeID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
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

		if req.UserID == "" || req.ShowtimeID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user_id and showtime_id are required",
			})
			return
		}

		updated, err := seat.UpdateSeatStatusIfCurrent(req.ShowtimeID, seatNumber, "LOCKED", "BOOKED")
		if err != nil {
			_ = audit.LogEvent("SYSTEM_ERROR", req.UserID, seatNumber, req.ShowtimeID, err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if !updated {
			c.JSON(http.StatusConflict, gin.H{
				"message": "seat is not in LOCKED state",
			})
			return
		}

		ok, err := ConfirmBooking(seatNumber, req.ShowtimeID)
		if err != nil {
			_, _ = seat.UpdateSeatStatusIfCurrent(req.ShowtimeID, seatNumber, "BOOKED", "LOCKED")
			_ = audit.LogEvent("SYSTEM_ERROR", req.UserID, seatNumber, req.ShowtimeID, err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if !ok {
			_, _ = seat.UpdateSeatStatusIfCurrent(req.ShowtimeID, seatNumber, "BOOKED", "LOCKED")

			c.JSON(http.StatusConflict, gin.H{
				"message": "booking is not in PENDING state",
			})
			return
		}

		key := "seat_lock:" + req.ShowtimeID + ":" + seatNumber
		_ = redisClient.Client.Del(c.Request.Context(), key).Err()

		_ = ws.BroadcastSeatEvent(hub, ws.SeatEvent{
			Event:      "seat_updated",
			ShowtimeID: req.ShowtimeID,
			SeatNumber: seatNumber,
			Status:     "BOOKED",
			UserID:     req.UserID,
		})

		_ = audit.LogEvent("BOOKING_SUCCESS", req.UserID, seatNumber, req.ShowtimeID, "booking confirmed successfully")

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
