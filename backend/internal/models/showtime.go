package models

import "time"

type Showtime struct {
	ID         string    `json:"id" bson:"_id"`
	MovieID    string    `json:"movie_id" bson:"movie_id"`
	MovieTitle string    `json:"movie_title" bson:"movie_title"`
	StartTime  time.Time `json:"start_time" bson:"start_time"`
	Theater    string    `json:"theater" bson:"theater"`
}
