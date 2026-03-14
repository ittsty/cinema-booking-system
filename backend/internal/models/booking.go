package models

import "time"

type BookingStatus string

const (
	PENDING BookingStatus = "PENDING"
	BOOKED  BookingStatus = "BOOKED"
)

type Booking struct {
	ID         string        `bson:"_id,omitempty"`
	UserID     string        `bson:"user_id"`
	SeatNumber string        `bson:"seat_number"`
	ShowtimeID string        `bson:"showtime_id"`
	Status     BookingStatus `bson:"status"`
	CreatedAt  time.Time     `bson:"created_at"`
}
