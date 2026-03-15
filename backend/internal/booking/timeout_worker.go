package booking

import (
	"context"
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
		expired, err := ExpireBooking(b.SeatNumber, b.ShowtimeID)
		if err != nil {
			log.Println("expire booking error:", err)
			_ = audit.LogEvent("SYSTEM_ERROR", b.UserID, b.SeatNumber, b.ShowtimeID, err.Error())
			continue
		}

		if !expired {
			log.Println("booking was not expired because it is no longer in PENDING state:", b.SeatNumber)
			continue
		}

		key := "seat_lock:" + b.ShowtimeID + ":" + b.SeatNumber
		_ = redisClient.Client.Del(context.Background(), key).Err()

		updated, err := seat.UpdateSeatStatusIfCurrent(b.ShowtimeID, b.SeatNumber, "LOCKED", "AVAILABLE")
		if err != nil {
			log.Println("update seat available error:", err)
			_ = audit.LogEvent("SYSTEM_ERROR", b.UserID, b.SeatNumber, b.ShowtimeID, err.Error())
			continue
		}

		if !updated {
			log.Println("seat was not updated because current status was not LOCKED:", b.SeatNumber)
			continue
		}

		_ = ws.BroadcastSeatEvent(hub, ws.SeatEvent{
			Event:      "seat_updated",
			ShowtimeID: b.ShowtimeID,
			SeatNumber: b.SeatNumber,
			Status:     "AVAILABLE",
			UserID:     b.UserID,
		})

		_ = audit.LogEvent("BOOKING_TIMEOUT", b.UserID, b.SeatNumber, b.ShowtimeID, "booking expired and seat released")
	}
}
