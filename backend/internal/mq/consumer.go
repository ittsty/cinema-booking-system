package mq

import (
	"encoding/json"
	"log"

	"cinema-booking/internal/audit"
)

func StartBookingSuccessConsumer() {
	msgs, err := Channel.Consume(
		QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("failed to consume queue:", err)
		return
	}

	go func() {
		for msg := range msgs {
			var event BookingSuccessEvent

			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Println("failed to unmarshal booking event:", err)
				continue
			}

			log.Println("booking success consumer received:", event)

			// mock notification / async logging
			_ = audit.LogEvent(
				"NOTIFICATION_SENT",
				event.UserID,
				event.SeatNumber,
				event.ShowtimeID,
				"mock notification sent via consumer",
			)
		}
	}()
}
