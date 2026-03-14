package booking

import (
	"context"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateBooking(booking models.Booking) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongo.DB.Collection("bookings")

	_, err := collection.InsertOne(ctx, booking)
	return err
}

func ConfirmBooking(seatNumber string, showtimeID string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongo.DB.Collection("bookings")

	_, err := collection.UpdateOne(
		ctx,
		bson.M{
			"seat_number": seatNumber,
			"showtime_id": showtimeID,
			"status":      models.PENDING,
		},
		bson.M{
			"$set": bson.M{
				"status":     models.BOOKED,
				"updated_at": time.Now(),
			},
		},
	)

	return err
}
