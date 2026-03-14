package main

import (
	"cinema-booking/internal/booking"
	"cinema-booking/internal/seat"
	"cinema-booking/pkg/mongo"
	"cinema-booking/pkg/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	mongo.Connect()
	redis.Connect()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	router.GET("/showtimes/:id/seats", seat.GetSeatMap)
	router.POST("/seats/:seat_number/lock", seat.LockSeatHandler)
	router.POST("/booking", booking.CreateBookingHandler)
	router.POST("/booking/:seat_number/confirm", booking.ConfirmBookingHandler)
	router.Run(":8080")
}
