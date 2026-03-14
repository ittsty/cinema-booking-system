package audit

import (
	"context"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/pkg/mongo"
)

func CreateLog(log models.AuditLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongo.DB.Collection("audit_logs")

	_, err := collection.InsertOne(ctx, log)
	return err
}

func LogEvent(event, userID, seatNumber, showtimeID, message string) error {
	log := models.AuditLog{
		Event:      event,
		UserID:     userID,
		SeatNumber: seatNumber,
		ShowtimeID: showtimeID,
		Message:    message,
		CreatedAt:  time.Now(),
	}

	return CreateLog(log)
}
