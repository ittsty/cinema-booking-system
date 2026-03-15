package seat

import (
	"context"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/pkg/mongo"
	redisClient "cinema-booking/pkg/redis"

	"go.mongodb.org/mongo-driver/bson"
)

func getSeatLockKey(showtimeID string, seatNumber string) string {
	return "seat_lock:" + showtimeID + ":" + seatNumber
}

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

		key := getSeatLockKey(showtimeID, seats[i].SeatNumber)
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

func LockSeat(showtimeID string, seatNumber string, userID string) (bool, error) {

	ctx := context.Background()

	key := getSeatLockKey(showtimeID, seatNumber)

	ok, err := redisClient.Client.SetNX(
		ctx,
		key,
		userID,
		5*time.Minute,
	).Result()

	return ok, err
}

func UnlockSeat(showtimeID string, seatNumber string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := getSeatLockKey(showtimeID, seatNumber)
	return redisClient.Client.Del(ctx, key).Err()
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

func UpdateSeatStatusIfCurrent(showtimeID string, seatNumber string, currentStatus string, newStatus string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongo.DB.Collection("seats")

	result, err := collection.UpdateOne(
		ctx,
		bson.M{
			"showtime_id": showtimeID,
			"seat_number": seatNumber,
			"status":      currentStatus,
		},
		bson.M{
			"$set": bson.M{
				"status": newStatus,
			},
		},
	)
	if err != nil {
		return false, err
	}

	return result.ModifiedCount == 1, nil
}

func GetSeatByShowtimeAndNumber(showtimeID string, seatNumber string) (*models.Seat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongo.DB.Collection("seats")

	var seat models.Seat
	err := collection.FindOne(ctx, bson.M{
		"showtime_id": showtimeID,
		"seat_number": seatNumber,
	}).Decode(&seat)
	if err != nil {
		return nil, err
	}

	key := getSeatLockKey(showtimeID, seatNumber)
	exists, err := redisClient.Client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if seat.Status != "BOOKED" && exists == 1 {
		seat.Status = "LOCKED"
	}

	return &seat, nil
}
