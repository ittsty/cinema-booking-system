package main

import (
	"cinema-booking/internal/admin"
	"cinema-booking/internal/audit"
	"cinema-booking/internal/auth"
	"cinema-booking/internal/booking"
	"cinema-booking/internal/mq"
	"cinema-booking/internal/seat"
	"cinema-booking/internal/ws"
	"cinema-booking/pkg/config"
	"cinema-booking/pkg/firebase"
	middleware "cinema-booking/pkg/middlewares"
	"cinema-booking/pkg/mongo"
	"cinema-booking/pkg/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	firebase.Init()

	mongo.Connect()
	redis.Connect()
	mq.Connect()

	hub := ws.NewHub()
	go hub.Run()

	booking.StartTimeoutWorker(hub)
	mq.StartBookingSuccessConsumer()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWS(hub, c)
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	router.GET("/showtimes", seat.GetShowtimesHandler)
	router.GET("/showtimes/:id/seats", seat.GetSeatMap)

	authGroup := router.Group("/")
	authGroup.Use(middleware.AuthRequired())
	{
		authGroup.GET("/me", auth.MeHandler)

		authGroup.POST("/seats/:seatNumber/lock", seat.LockSeatHandler(hub))
		authGroup.POST("/seats/:seatNumber/unlock", seat.UnlockSeatHandler(hub))

		authGroup.POST("/booking", booking.CreateBookingHandler)
		authGroup.POST("/booking/:seat_number/confirm", booking.ConfirmBookingHandler(hub))
	}

	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AuthRequired(), middleware.AdminOnly())
	{
		adminGroup.GET("/bookings", admin.GetBookingsHandler)
		adminGroup.GET("/logs", audit.GetLogsHandler)
	}
	router.Run(":" + config.App.Port)
}
