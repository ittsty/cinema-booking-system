package main

import (
	"cinema-booking/internal/admin"
	"cinema-booking/internal/audit"
	"cinema-booking/internal/booking"
	"cinema-booking/internal/mq"
	"cinema-booking/internal/seat"
	"cinema-booking/internal/ws"
	"cinema-booking/pkg/config"
	middleware "cinema-booking/pkg/middlewares"
	"cinema-booking/pkg/mongo"
	"cinema-booking/pkg/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()

	mongo.Connect()
	redis.Connect()
	mq.Connect()

	hub := ws.NewHub()
	go hub.Run()

	booking.StartTimeoutWorker(hub)
	mq.StartBookingSuccessConsumer()

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

	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AdminOnly())
	{
		adminGroup.GET("/bookings", admin.GetBookingsHandler)
		adminGroup.GET("/logs", audit.GetLogsHandler)
	}
	router.Run(":" + config.App.Port)
}
