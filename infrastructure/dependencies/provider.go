package dependencies

import (
	"sync"
	"time"
	"worker-nicepay/application/services"
	"worker-nicepay/infrastructure/configuration"
	"worker-nicepay/infrastructure/database"
	"worker-nicepay/infrastructure/database/connectors"
	"worker-nicepay/infrastructure/database/repositories"
	"worker-nicepay/infrastructure/gateway/nicepay"
	"worker-nicepay/infrastructure/publishers"
	"worker-nicepay/infrastructure/service"

	"github.com/google/wire"
)

// singleton
var gatewayOnce sync.Once
var transactionServiceOnce sync.Once
var publisherOnce sync.Once
var yugabyteClientOnce sync.Once
var masterDataRepoOnce sync.Once
var paymentRepoOnce sync.Once
var xenditRepoOnce sync.Once

// singleton instance
var nicepayGatewayInstance *nicepay.NicepayGateway
var publisherInstance *publishers.PublisherLog
var yugabyteClientInstance *connectors.YugabyteConnector
var masterDataRepoInstance *repositories.MasterDataRepositoryYugabyteDB
var paymentRepoInstance *repositories.PaymentRepositoryYugabyteDB
var xenditRepoInstance *repositories.XenditRepositoryYugabyteDB
var currenciesRepoInstance *repositories.CurrenciesRepository
var countriesRepoInstance *repositories.CountriesRepository
var merchantsRepoInstance *repositories.MerchantsRepository
var paymentMethodsRepoInstance *repositories.PaymentMethodsRepository
var NicepaytransactionServiceInstance *service.NicePayTransactionService

var ProviderSet wire.ProviderSet = wire.NewSet(
	ProvideNicepayGateway,
	ProvideTransactionService,
	ProvideYugabyteClient,
	ProvideMasterDataRepository,
	ProvidePaymentRepository,
	ProvideXenditRepository,
	ProvideCurrenciesRepository,
	ProvideCountriesRepository,
	ProvideMerchantsRepository,
	ProvidePaymentMethodsRepository,
	ProvidePublisher,
	wire.Bind(new(services.PaymentGateway), new(*nicepay.NicepayGateway)),
	wire.Bind(new(services.TransactionService), new(*service.NicePayTransactionService)),
	wire.Bind(new(services.Publisher), new(*publishers.PublisherLog)),
)

func ProvideNicepayGateway() *nicepay.NicepayGateway {
	gatewayOnce.Do(func() {
		timeout := time.Duration(configuration.AppConfig.XenditTimeout) * time.Millisecond
		nicepayGatewayInstance = nicepay.NewNicepayGateway(configuration.AppConfig.XenditAPIURL, configuration.AppConfig.XenditAPIKey, timeout)
	})
	return nicepayGatewayInstance
}

func ProvideTransactionService() *service.NicePayTransactionService {
	transactionServiceOnce.Do(func() {
		// masterRepo := ProvideMasterDataRepository()
		paymentRepo := ProvidePaymentRepository()
		currencyRepo := ProvideCurrenciesRepository()
		countryRepo := ProvideCountriesRepository()
		merchantRepo := ProvideMerchantsRepository()
		paymentMethodRepo := ProvidePaymentMethodsRepository()
		gateway := ProvideNicepayGateway()
		db := ProvideYugabyteClient().GetDB()
		NicepaytransactionServiceInstance = service.NewNicePayTransactionService(db, paymentRepo, currencyRepo, countryRepo, paymentMethodRepo, merchantRepo, gateway)
	})
	return NicepaytransactionServiceInstance
}

func ProvideYugabyteClient() *connectors.YugabyteConnector {
	yugabyteClientOnce.Do(func() {
		yugabyteClientInstance = connectors.NewYugabyteConnector(database.YugabyteDBClient)
	})
	return yugabyteClientInstance
}

func ProvideMasterDataRepository() *repositories.MasterDataRepositoryYugabyteDB {
	masterDataRepoOnce.Do(func() {
		masterDataRepoInstance = repositories.NewMasterDataRepositoryYugabyteDB()
	})
	return masterDataRepoInstance
}

func ProvidePaymentRepository() *repositories.PaymentRepositoryYugabyteDB {
	paymentRepoOnce.Do(func() {
		paymentRepoInstance = repositories.NewPaymentRepositoryYugabyteDB()
	})
	return paymentRepoInstance
}

func ProvideXenditRepository() *repositories.XenditRepositoryYugabyteDB {
	xenditRepoOnce.Do(func() {
		xenditRepoInstance = repositories.NewXenditRepositoryYugabyteDB()
	})
	return xenditRepoInstance
}

func ProvidePublisher() *publishers.PublisherLog {
	publisherOnce.Do(func() {
		publisherInstance = publishers.NewPublisherLog()
	})
	return publisherInstance
}

func ProvideCurrenciesRepository() *repositories.CurrenciesRepository {
	if currenciesRepoInstance == nil {
		currenciesRepoInstance = repositories.NewCurrenciesRepository()
	}
	return currenciesRepoInstance
}

func ProvideCountriesRepository() *repositories.CountriesRepository {
	if countriesRepoInstance == nil {
		countriesRepoInstance = repositories.NewCountriesRepository()
	}
	return countriesRepoInstance
}

func ProvideMerchantsRepository() *repositories.MerchantsRepository {
	if merchantsRepoInstance == nil {
		merchantsRepoInstance = repositories.NewMerchantsRepository()
	}
	return merchantsRepoInstance
}

func ProvidePaymentMethodsRepository() *repositories.PaymentMethodsRepository {
	if paymentMethodsRepoInstance == nil {
		paymentMethodsRepoInstance = repositories.NewPaymentMethodsRepository()
	}
	return paymentMethodsRepoInstance
}
