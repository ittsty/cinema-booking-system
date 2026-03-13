package models

type SeatStatus string

const (
	SeatAvailable SeatStatus = "AVAILABLE"
	SeatLocked    SeatStatus = "LOCKED"
	SeatBooked    SeatStatus = "BOOKED"
)

type Seat struct {
	ID         string     `bson:"_id,omitempty" json:"id"`
	ShowtimeID string     `bson:"showtime_id" json:"showtime_id"`
	SeatNumber string     `bson:"seat_number" json:"seat_number"`
	Status     SeatStatus `bson:"status" json:"status"`
}