package services

import (
	"context"
	"worker-nicepay/infrastructure/gateway/nicepay"
)

type PaymentGateway interface {
	RequestPaymentLink(ctx context.Context, req nicepay.RequestPaymentLinkDTO, url string) (nicepay.ResponsePaymentLinkDTO, error)
}
