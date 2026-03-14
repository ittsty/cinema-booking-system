package seat

import (
	"context"
	"fmt"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/internal/ws"
	"cinema-booking/pkg/mongo"
	"cinema-booking/pkg/redis"
	redisClient "cinema-booking/pkg/redis"

	"go.mongodb.org/mongo-driver/bson"
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
	if err := cursor.All(ctx, &seats); err != nil {
		return nil, err
	}

	for i := range seats {
		if seats[i].Status == "BOOKED" {
			continue
		}

		key := "seat_lock:" + seats[i].SeatNumber
		exists, err := redisClient.Client.Exists(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		if exists == 1 {
			seats[i].Status = "LOCKED"
		}
	}
	return seats, nil
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

func UpdateSeatStatus(showtimeID string, seatNumber string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongo.DB.Collection("seats")

	_, err := collection.UpdateOne(
		ctx,
		bson.M{
			"showtime_id": showtimeID,
			"seat_number": seatNumber,
		},
		bson.M{
			"$set": bson.M{
				"status": status,
			},
		},
	)

	return err
}
