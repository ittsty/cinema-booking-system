package seat

import (
	"context"
	"fmt"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/internal/ws"
	"cinema-booking/pkg/mongo"
	"cinema-booking/pkg/redis"
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

func LockSeat(seatNumber string, userID string, hub *ws.Hub) (bool, error) {

	ctx := context.Background()

	key := "seat_lock:" + seatNumber

	ok, err := redis.Client.SetNX(
		ctx,
		key,
		userID,
		5*time.Minute,
	).Result()
	if ok {

		msg := fmt.Sprintf(`{
 		 "event":"seat_locked",
  		 "seat":"%s"
 		}`, seatNumber)

		hub.Broadcast <- []byte(msg)
	}
	return ok, err
}
