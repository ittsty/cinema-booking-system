package models

import "time"

type AuditLog struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	Event      string    `bson:"event" json:"event"`
	UserID     string    `bson:"user_id,omitempty" json:"user_id,omitempty"`
	SeatNumber string    `bson:"seat_number,omitempty" json:"seat_number,omitempty"`
	ShowtimeID string    `bson:"showtime_id,omitempty" json:"showtime_id,omitempty"`
	Message    string    `bson:"message" json:"message"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}
