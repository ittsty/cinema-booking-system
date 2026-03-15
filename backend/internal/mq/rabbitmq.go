package mq

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

const QueueName = "booking.success"

func Connect() {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@rabbitmq:5672/"
	}

	var err error
	Conn, err = amqp.Dial(url)
	if err != nil {
		log.Fatal("failed to connect rabbitmq:", err)
	}

	Channel, err = Conn.Channel()
	if err != nil {
		log.Fatal("failed to open rabbitmq channel:", err)
	}

	_, err = Channel.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("failed to declare queue:", err)
	}
}
