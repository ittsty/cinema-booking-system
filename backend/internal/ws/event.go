package ws

type SeatEvent struct {
	Event      string `json:"event"`
	ShowtimeID string `json:"showtime_id"`
	SeatNumber string `json:"seat_number"`
	Status     string `json:"status"`
	UserID     string `json:"user_id,omitempty"`
}
