package seat

import (
	"context"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/pkg/mongo"
)

func GetSeats(showtimeID string) ([]models.Seat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongo.DB.Collection("seats")

	cursor, err := collection.Find(ctx, map[string]string{
		"showtime_id": showtimeID,
	})

	if err != nil {
		return nil, err
	}

	var seats []models.Seat

	err = cursor.All(ctx, &seats)

	return seats, err
}
