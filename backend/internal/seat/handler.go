package seat

import (
	"net/http"

	"cinema-booking/internal/audit"
	"cinema-booking/internal/ws"
	redisClient "cinema-booking/pkg/redis"

	"github.com/gin-gonic/gin"
)

func GetShowtimesHandler(c *gin.Context) {
	showtimes, err := GetShowtimes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, showtimes)
}

func GetSeatMap(c *gin.Context) {

	showtimeID := c.Param("id")
	seats, err := GetSeats(showtimeID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, seats)
}

func LockSeatHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		seatNumber := c.Param("seatNumber")
		userID := c.GetString("user_id")
		showtimeID := c.Query("showtime_id")

		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		if showtimeID == "" || seatNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "showtime_id and seatNumber are required",
			})
			return
		}

		ok, err := LockSeat(showtimeID, seatNumber, userID)
		if err != nil {
			_ = audit.LogEvent("SYSTEM_ERROR", userID, seatNumber, showtimeID, err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if !ok {
			_ = audit.LogEvent("SEAT_LOCK_FAILED", userID, seatNumber, showtimeID, "seat already locked")

			c.JSON(http.StatusConflict, gin.H{
				"message": "seat already locked",
			})
			return
		}
		updated, err := UpdateSeatStatusIfCurrent(showtimeID, seatNumber, "AVAILABLE", "LOCKED")
		if err != nil {
			_ = UnlockSeat(showtimeID, seatNumber)
			_ = audit.LogEvent("SYSTEM_ERROR", userID, seatNumber, showtimeID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if !updated {
			_ = UnlockSeat(showtimeID, seatNumber)
			c.JSON(http.StatusConflict, gin.H{
				"message": "seat is not in AVAILABLE state",
			})
			return
		}
		_ = ws.BroadcastSeatEvent(hub, ws.SeatEvent{
			Event:      "seat_updated",
			ShowtimeID: showtimeID,
			SeatNumber: seatNumber,
			Status:     "LOCKED",
			UserID:     userID,
		})
		_ = audit.LogEvent("SEAT_LOCKED", userID, seatNumber, showtimeID, "seat locked successfully")

		c.JSON(http.StatusOK, gin.H{
			"message": "seat locked",
		})
	}
}

func UnlockSeatHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		seatNumber := c.Param("seatNumber")
		showtimeID := c.Query("showtime_id")
		userID := c.GetString("user_id")

		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		if showtimeID == "" || seatNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "showtime_id and seatNumber are required",
			})
			return
		}

		key := "seat_lock:" + showtimeID + ":" + seatNumber
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

		if err := UnlockSeat(showtimeID, seatNumber); err != nil {
			_ = audit.LogEvent("SYSTEM_ERROR", userID, seatNumber, showtimeID, err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		updated, err := UpdateSeatStatusIfCurrent(showtimeID, seatNumber, "LOCKED", "AVAILABLE")
		if err != nil {
			_ = audit.LogEvent("SYSTEM_ERROR", userID, seatNumber, showtimeID, err.Error())

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

		_ = ws.BroadcastSeatEvent(hub, ws.SeatEvent{
			Event:      "seat_updated",
			ShowtimeID: showtimeID,
			SeatNumber: seatNumber,
			Status:     "AVAILABLE",
			UserID:     userID,
		})

		_ = audit.LogEvent("SEAT_RELEASED", userID, seatNumber, showtimeID, "seat unlocked successfully")

		c.JSON(http.StatusOK, gin.H{
			"message": "seat unlocked",
		})
	}
}
