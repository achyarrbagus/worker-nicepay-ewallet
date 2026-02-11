package publishers

import (
	"context"
	"encoding/json"
	"log"
	"worker-nicepay/application/messages"
	"worker-nicepay/application/services"
)

type PublisherLog struct{}

func NewPublisherLog() *PublisherLog {
	return &PublisherLog{}
}

func (p *PublisherLog) Publish(ctx context.Context, message services.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}
	switch message.GetMessageName() {
	case messages.PaymentCreatedMessageName:
		log.Printf("Publishing message: %s - %s", message.GetMessageName(), string(payload))
	default:
		log.Printf("Publishing unknown message: %T", message)
	}
	return nil
}
