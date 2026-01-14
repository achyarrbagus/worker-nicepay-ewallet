package services

import (
	"context"
	"fmt"

	"payment-airpay/domain/entities"
)

type CreatePaymentService struct {
	Gateway PaymentGateway
	TxSvc   TransactionService
}

func NewCreatePaymentService(g PaymentGateway, t TransactionService) *CreatePaymentService {
	return &CreatePaymentService{Gateway: g, TxSvc: t}
}

func (s *CreatePaymentService) Execute(ctx context.Context, payload map[string]interface{}) (entities.Payment, error) {
	res, err := s.Gateway.Create(ctx, payload)
	if err != nil {
		return entities.Payment{}, err
	}
	if err := s.TxSvc.Save(ctx, res, payload); err != nil {
		return entities.Payment{}, fmt.Errorf("failed to persist payment: %w", err)
	}
	return res, nil
}
