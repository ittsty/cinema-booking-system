package models

import "time"

type Showtime struct {
	ID        string    `bson:"_id,omitempty"`
	MovieID   string    `bson:"movie_id"`
	StartTime time.Time `bson:"start_time"`
}