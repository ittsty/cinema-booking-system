package main

import (
	"cinema-booking/internal/audit"
	"cinema-booking/internal/booking"
	"cinema-booking/internal/seat"
	"cinema-booking/internal/ws"
	"cinema-booking/pkg/mongo"
	"cinema-booking/pkg/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	mongo.Connect()
	redis.Connect()

	hub := ws.NewHub()
	go hub.Run()

	booking.StartTimeoutWorker(hub)

	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWS(hub, c)
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	router.GET("/showtimes/:id/seats", seat.GetSeatMap)
	router.POST("/seats/:seatNumber/lock", seat.LockSeatHandler(hub))

	router.POST("/booking", booking.CreateBookingHandler)
	router.POST("/booking/:seat_number/confirm", booking.ConfirmBookingHandler(hub))

	router.GET("/admin/logs", audit.GetLogsHandler)
	router.Run(":8080")
}
