package booking

import (
	"context"
	"fmt"
	"log"
	"time"

	"cinema-booking/internal/audit"
	"cinema-booking/internal/seat"
	"cinema-booking/internal/ws"
	redisClient "cinema-booking/pkg/redis"
)

func StartTimeoutWorker(hub *ws.Hub) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			processExpiredBookings(hub)
		}
	}()
}

func processExpiredBookings(hub *ws.Hub) {
	bookings, err := FindExpiredPendingBookings(5 * time.Minute)
	if err != nil {
		log.Println("timeout worker find expired bookings error:", err)
		return
	}

	for _, b := range bookings {
		if err := ExpireBooking(b.SeatNumber, b.ShowtimeID); err != nil {
			log.Println("expire booking error:", err)
			_ = audit.LogEvent("SYSTEM_ERROR", b.UserID, b.SeatNumber, b.ShowtimeID, err.Error())
			continue
		}

		key := "seat_lock:" + b.SeatNumber
		_ = redisClient.Client.Del(context.Background(), key).Err()

		if err := seat.UpdateSeatStatus(b.ShowtimeID, b.SeatNumber, "AVAILABLE"); err != nil {
			log.Println("update seat available error:", err)
			_ = audit.LogEvent("SYSTEM_ERROR", b.UserID, b.SeatNumber, b.ShowtimeID, err.Error())
			continue
		}

		msg := fmt.Sprintf(`{
			"event":"seat_released",
			"seat":"%s"
		}`, b.SeatNumber)

		hub.Broadcast <- []byte(msg)

		_ = audit.LogEvent("BOOKING_TIMEOUT", b.UserID, b.SeatNumber, b.ShowtimeID, "booking expired and seat released")
	}
}
