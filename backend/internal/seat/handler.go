package seat

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func LockSeatHandler(c *gin.Context) {

	seatNumber := c.Param("seat_number")
	userID := c.Query("user_id")
	ok, err := LockSeat(seatNumber, userID)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !ok {
		c.JSON(409, gin.H{
			"message": "seat already locked",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "seat locked successfully",
	})
}
