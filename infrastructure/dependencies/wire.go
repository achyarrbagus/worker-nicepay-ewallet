//go:build wireinject
// +build wireinject

package dependencies

import (
	"worker-nicepay/application/services"
	"worker-nicepay/infrastructure/gateway/nicepay"
	"worker-nicepay/infrastructure/publishers"
	"worker-nicepay/infrastructure/service"

	"github.com/google/wire"
)

func WireCreatePaymentService() *services.CreatePaymentService {
	panic(wire.Build(ProviderSet, services.NewCreatePaymentService))
}

func WireNicepayGateway() *nicepay.NicepayGateway {
	panic(wire.Build(ProviderSet))
}

func WireNicepayTransactionService() *service.NicePayTransactionService {
	panic(wire.Build(ProviderSet))
}

func WirePublisher() *publishers.PublisherLog {
	panic(wire.Build(ProviderSet))
}
