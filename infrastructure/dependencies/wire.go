//go:build wireinject
// +build wireinject

package dependencies

import (
	"payment-airpay/application/services"
	"payment-airpay/infrastructure/gateway/xendit"
	"payment-airpay/infrastructure/publishers"
	"payment-airpay/infrastructure/service"

	"github.com/google/wire"
)

func WireCreatePaymentService() *services.CreatePaymentService {
	panic(wire.Build(ProviderSet, services.NewCreatePaymentService))
}

func WireXenditGateway() *xendit.XenditGateway {
	panic(wire.Build(ProviderSet))
}

func WireTransactionService() *service.PaymentXendit {
	panic(wire.Build(ProviderSet))
}

func WirePublisher() *publishers.PublisherLog {
	panic(wire.Build(ProviderSet))
}
