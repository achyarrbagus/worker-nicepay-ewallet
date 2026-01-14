package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"payment-airpay/application/events"
	"payment-airpay/application/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQQueue implements services.Queue and publishes supported tasks/events to RabbitMQ
type RabbitMQQueue struct {
	ch *amqp.Channel
}

// NewRabbitMQQueue constructs a RabbitMQQueue using the global RabbitChan
func NewRabbitMQQueue(channel *amqp.Channel) *RabbitMQQueue {
	return &RabbitMQQueue{ch: channel}
}

// Enqueue publishes a task to RabbitMQ. This implementation only supports
// publishing PaymentCreatedEvent. It uses an exchange named after the event
// (fanout) and publishes the JSON payload with persistent delivery mode.
func (r *RabbitMQQueue) Enqueue(ctx context.Context, event services.Event) error {
	if r == nil || r.ch == nil {
		return errors.New("rabbitmq channel is not initialized; call InitializeRabbitMQ first")
	}

	// Only allow the PaymentCreatedEvent to be published via RabbitMQQueue
	allowedEventName := (events.PaymentCreatedEvent{}).GetEventName()
	eventName := event.GetEventName()
	if eventName != allowedEventName {
		return fmt.Errorf("unsupported event for RabbitMQQueue: %s (only %s is supported)", eventName, allowedEventName)
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	// Publish to the exchange with empty routing key (fanout)
	pub := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         payload,
		Timestamp:    time.Now(),
	}

	exchangeName := strings.ReplaceAll(eventName, "-", ".")
	if err := r.ch.PublishWithContext(ctx,
		exchangeName, // exchange
		"",           // routing key (ignored for fanout)
		false,        // mandatory
		false,        // immediate
		pub,
	); err != nil {
		return fmt.Errorf("failed to publish message to %s: %w", eventName, err)
	}

	return nil
}
