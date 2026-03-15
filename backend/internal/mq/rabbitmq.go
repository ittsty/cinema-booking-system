package mq

import (
	"log"

	"cinema-booking/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

const QueueName = "booking.success"

func Connect() {
	var err error

	Conn, err = amqp.Dial(config.App.RabbitMQURL)
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

	log.Println("RabbitMQ connected")
}
