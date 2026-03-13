package models

import "time"

type Booking struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	UserID     string    `bson:"user_id"`
	ShowtimeID string    `bson:"showtime_id"`
	SeatID     string    `bson:"seat_id"`
	Status     string    `bson:"status"`
	CreatedAt  time.Time `bson:"created_at"`
}