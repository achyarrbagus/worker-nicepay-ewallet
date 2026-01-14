package dependencies

import (
	"payment-airpay/application/services"
	"payment-airpay/infrastructure/configuration"
	"payment-airpay/infrastructure/database"
	"payment-airpay/infrastructure/database/connectors"
	"payment-airpay/infrastructure/database/repositories"
	"payment-airpay/infrastructure/gateway/xendit"
	"payment-airpay/infrastructure/publishers"
	"payment-airpay/infrastructure/service"
	"sync"
	"time"

	"github.com/google/wire"
)

// singleton
var gatewayOnce sync.Once
var transactionServiceOnce sync.Once
var publisherOnce sync.Once
var yugabyteClientOnce sync.Once
var paymentXenditRepoOnce sync.Once

// singleton instance
var xenditGatewayInstance *xendit.XenditGateway
var transactionServiceInstance *service.PaymentXendit
var publisherInstance *publishers.PublisherLog
var yugabyteClientInstance *connectors.YugabyteConnector
var paymentXenditRepoInstance *repositories.PaymentXenditRepositoryYugabyteDB

var ProviderSet wire.ProviderSet = wire.NewSet(
	ProvideXenditGateway,
	ProvideTransactionService,
	ProvideYugabyteClient,
	ProvidePaymentXenditRepository,
	ProvidePublisher,
	wire.Bind(new(services.PaymentGateway), new(*xendit.XenditGateway)),
	wire.Bind(new(services.TransactionService), new(*service.PaymentXendit)),
	wire.Bind(new(services.Publisher), new(*publishers.PublisherLog)),
)

func ProvideXenditGateway() *xendit.XenditGateway {
	gatewayOnce.Do(func() {
		timeout := time.Duration(configuration.AppConfig.XenditTimeout) * time.Millisecond
		xenditGatewayInstance = xendit.NewXenditGateway(configuration.AppConfig.XenditAPIURL, configuration.AppConfig.XenditAPIKey, timeout)
	})
	return xenditGatewayInstance
}

func ProvideTransactionService() *service.PaymentXendit {
	transactionServiceOnce.Do(func() {
		repo := ProvidePaymentXenditRepository()
		db := ProvideYugabyteClient()
		transactionServiceInstance = service.NewPaymentXendit(repo, db)
	})
	return transactionServiceInstance
}

func ProvideYugabyteClient() *connectors.YugabyteConnector {
	yugabyteClientOnce.Do(func() {
		yugabyteClientInstance = connectors.NewYugabyteConnector(database.YugabyteDBClient)
	})
	return yugabyteClientInstance
}

func ProvidePaymentXenditRepository() *repositories.PaymentXenditRepositoryYugabyteDB {
	paymentXenditRepoOnce.Do(func() {
		paymentXenditRepoInstance = repositories.NewPaymentXenditRepositoryYugabyteDB(ProvideYugabyteClient())
	})
	return paymentXenditRepoInstance
}

func ProvidePublisher() *publishers.PublisherLog {
	publisherOnce.Do(func() {
		publisherInstance = publishers.NewPublisherLog()
	})
	return publisherInstance
}
