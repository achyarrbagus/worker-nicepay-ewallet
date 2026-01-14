package services

import (
	"context"

	"payment-airpay/domain/entities"
)

type TransactionService interface {
	Save(ctx context.Context, payment entities.Payment, payload map[string]interface{}) error
}
