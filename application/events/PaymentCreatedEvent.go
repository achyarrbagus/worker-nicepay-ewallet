package events

import "time"

type PaymentCreatedEvent struct {
	Timestamp     time.Time              `json:"timestamp"`
	TransactionID string                 `json:"transaction_id"`
	Message       string                 `json:"message"`
	Payload       map[string]interface{} `json:"payload"`
}

func (e PaymentCreatedEvent) GetEventName() string {
	return PaymentCreatedEventName
}
