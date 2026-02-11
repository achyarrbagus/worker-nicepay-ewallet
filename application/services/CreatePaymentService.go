package services

import (
	"context"
	"fmt"

	"worker-nicepay/application/dto"
	"worker-nicepay/domain/entities"
)

type CreatePaymentService struct {
	Gateway PaymentGateway
	TxSvc   TransactionService
}

func NewCreatePaymentService(g PaymentGateway, t TransactionService) *CreatePaymentService {
	return &CreatePaymentService{Gateway: g, TxSvc: t}
}

func (s *CreatePaymentService) Execute(ctx context.Context, req dto.CreatePaymentRequest, incoming entities.Incoming) (string, entities.Payment, error) {

	payementLinkUrl, payment, err := s.TxSvc.Save(ctx, req, incoming)
	if err != nil {
		return "", entities.Payment{}, fmt.Errorf("failed to persist payment: %w", err)
	}
	return payementLinkUrl, payment, nil
}
