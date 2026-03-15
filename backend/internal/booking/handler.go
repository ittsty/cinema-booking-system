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
	SeatNumber string `json:"seat_number"`
	ShowtimeID string `json:"showtime_id"`
}

type ConfirmBookingRequest struct {
	ShowtimeID string `json:"showtime_id"`
}

func CreateBookingHandler(c *gin.Context) {
	var req BookingRequest

	userID := c.GetString("user_id")
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = audit.LogEvent("SYSTEM_ERROR", userID, req.SeatNumber, req.ShowtimeID, err.Error())

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	if req.SeatNumber == "" || req.ShowtimeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "seat_number and showtime_id are required",
		})
		return
	}

	seatData, err := seat.GetSeatByShowtimeAndNumber(req.ShowtimeID, req.SeatNumber)
	if err != nil {
		_ = audit.LogEvent("SYSTEM_ERROR", userID, req.SeatNumber, req.ShowtimeID, err.Error())
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

	key := "seat_lock:" + req.ShowtimeID + ":" + req.SeatNumber
	lockedBy, err := redisClient.Client.Get(c.Request.Context(), key).Result()
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "seat is not locked",
		})
		return
	}

	if lockedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "you do not own this seat lock",
		})
		return
	}

	booking := models.Booking{
		UserID:     userID,
		SeatNumber: req.SeatNumber,
		ShowtimeID: req.ShowtimeID,
		Status:     models.PENDING,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	exists, err := HasPendingBooking(req.SeatNumber, req.ShowtimeID)
	if err != nil {
		_ = audit.LogEvent("SYSTEM_ERROR", userID, req.SeatNumber, req.ShowtimeID, err.Error())
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
		_ = audit.LogEvent("SYSTEM_ERROR", userID, req.SeatNumber, req.ShowtimeID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	_ = audit.LogEvent("BOOKING_CREATED", userID, req.SeatNumber, req.ShowtimeID, "booking created successfully")

	c.JSON(http.StatusCreated, gin.H{
		"message": "booking created",
		"data":    booking,
	})
}

func ConfirmBookingHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		seatNumber := c.Param("seat_number")
		userID := c.GetString("user_id")
		var req ConfirmBookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			_ = audit.LogEvent("SYSTEM_ERROR", userID, seatNumber, req.ShowtimeID, err.Error())

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		if req.ShowtimeID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "showtime_id is required",
			})
			return
		}
		key := "seat_lock:" + req.ShowtimeID + ":" + seatNumber
		lockedBy, err := redisClient.Client.Get(c.Request.Context(), key).Result()
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{
				"message": "seat is not locked",
			})
			return
		}

		if lockedBy != userID {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "you do not own this seat lock",
			})
			return
		}
		updated, err := seat.UpdateSeatStatusIfCurrent(req.ShowtimeID, seatNumber, "LOCKED", "BOOKED")
		if err != nil {
			_ = audit.LogEvent("SYSTEM_ERROR", userID, seatNumber, req.ShowtimeID, err.Error())

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
			_ = audit.LogEvent("SYSTEM_ERROR", userID, seatNumber, req.ShowtimeID, err.Error())

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

		_ = redisClient.Client.Del(c.Request.Context(), key).Err()

		_ = ws.BroadcastSeatEvent(hub, ws.SeatEvent{
			Event:      "seat_updated",
			ShowtimeID: req.ShowtimeID,
			SeatNumber: seatNumber,
			Status:     "BOOKED",
			UserID:     userID,
		})

		_ = audit.LogEvent("BOOKING_SUCCESS", userID, seatNumber, req.ShowtimeID, "booking confirmed successfully")

		event := mq.BookingSuccessEvent{
			UserID:     userID,
			SeatNumber: seatNumber,
			ShowtimeID: req.ShowtimeID,
			Status:     "BOOKED",
			CreatedAt:  time.Now(),
		}

		if err := mq.PublishBookingSuccess(event); err != nil {
			_ = audit.LogEvent("SYSTEM_ERROR", userID, seatNumber, req.ShowtimeID, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "booking confirmed",
		})
	}

}
