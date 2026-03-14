package booking

import (
	"context"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/pkg/mongo"
)

func CreateBooking(booking models.Booking) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongo.DB.Collection("bookings")

	_, err := collection.InsertOne(ctx, booking)

	return err
}

func ConfirmBooking(seatNumber string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongo.DB.Collection("bookings")

	_, err := collection.UpdateOne(
		ctx,
		map[string]string{
			"seat_number": seatNumber,
		},
		map[string]interface{}{
			"$set": map[string]string{
				"status": "BOOKED",
			},
		},
	)

	return err
}
