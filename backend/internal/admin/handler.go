package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBookingsHandler(c *gin.Context) {
	filter := BookingFilter{
		UserID:     c.GetString("user_id"),
		ShowtimeID: c.Query("showtime_id"),
		Status:     c.Query("status"),
	}

	bookings, err := GetBookings(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  bookings,
		"count": len(bookings),
	})
}
