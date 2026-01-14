package services

import (
	"context"
	"payment-airpay/domain/entities"
)

type PaymentGateway interface {
	Create(ctx context.Context, payload map[string]interface{}) (entities.Payment, error)
}
