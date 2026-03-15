package admin

import (
	"context"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookingFilter struct {
	UserID     string
	ShowtimeID string
	Status     string
}

func GetBookings(filter BookingFilter) ([]models.Booking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongo.DB.Collection("bookings")

	query := bson.M{}

	if filter.UserID != "" {
		query["user_id"] = filter.UserID
	}

	if filter.ShowtimeID != "" {
		query["showtime_id"] = filter.ShowtimeID
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	opts := options.Find().SetSort(bson.D{
		{Key: "created_at", Value: -1},
	})

	cursor, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	var bookings []models.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}
