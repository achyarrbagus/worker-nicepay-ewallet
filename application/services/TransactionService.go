package services

import (
	"context"

	"worker-nicepay/application/dto"
	"worker-nicepay/domain/entities"
)

type TransactionService interface {
	Save(ctx context.Context, param dto.CreatePaymentRequest, incoming entities.Incoming) (string, entities.Payment, error)
}
