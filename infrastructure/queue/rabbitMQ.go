package queue

import (
	"log"
	"payment-airpay/application/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp.Connection
var RabbitChan *amqp.Channel

func InitializeRabbitMQ() {
	var err error
	RabbitConn, err = amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ: ", err)
	}

	RabbitChan, err = RabbitConn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel: ", err)
	}

	// Declare a Exchange to ensure it exists
	exchangeName := (events.PaymentCreatedEvent{}).GetEventName()
	err = RabbitChan.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare exchange: ", err)
	}

	// Declare queue to ensure it exists
	queueName := exchangeName
	_, err = RabbitChan.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare queue: ", err)
	}

	// Bind the queue to the exchange so fanout routes messages to this queue
	if err := RabbitChan.QueueBind(
		queueName,    // queue name
		"",           // routing key (ignored for fanout)
		exchangeName, // exchange
		false,        // no-wait
		nil,          // args
	); err != nil {
		log.Fatal("Failed to bind queue to exchange: ", err)
	}
}

func CloseRabbitMQ() {
	if RabbitConn != nil {
		RabbitConn.Close()
	}
	if RabbitChan != nil {
		RabbitChan.Close()
	}
}
