package models

import "time"

type BookingStatus string

const (
	PENDING        BookingStatus = "PENDING"
	BOOKED         BookingStatus = "BOOKED"
	BookingExpired BookingStatus = "EXPIRED"
)

type Booking struct {
	ID         string        `bson:"_id,omitempty" json:"id"`
	UserID     string        `bson:"user_id" json:"user_id"`
	SeatNumber string        `bson:"seat_number" json:"seat_number"`
	ShowtimeID string        `bson:"showtime_id" json:"showtime_id"`
	Status     BookingStatus `bson:"status" json:"status"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at" json:"updated_at"`
}
