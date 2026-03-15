package mq

import "time"

type BookingSuccessEvent struct {
	UserID     string    `json:"user_id"`
	SeatNumber string    `json:"seat_number"`
	ShowtimeID string    `json:"showtime_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
