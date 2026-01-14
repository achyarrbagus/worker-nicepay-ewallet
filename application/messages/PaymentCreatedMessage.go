package messages

import "time"

type PaymentCreatedMessage struct {
	Timestamp     time.Time `json:"timestamp"`
	TransactionID string    `json:"transaction_id"`
	Message       string    `json:"message"`
}

func (m PaymentCreatedMessage) GetMessageName() string {
	return PaymentCreatedMessageName
}
