package service

import (
	"context"
	"time"

	"worker-nicepay/application/dto"
	"worker-nicepay/application/services"
	"worker-nicepay/domain/entities"
	"worker-nicepay/infrastructure/configuration"
	constant "worker-nicepay/infrastructure/const"
	"worker-nicepay/infrastructure/database/models"
	"worker-nicepay/infrastructure/database/repositories"
	"worker-nicepay/infrastructure/gateway/nicepay"

	"gorm.io/gorm"
)

type NicePayTransactionService struct {
	db                *gorm.DB
	TransactionRepo   *repositories.PaymentRepositoryYugabyteDB
	CurrencyRepo      *repositories.CurrenciesRepository
	CountryRepo       *repositories.CountriesRepository
	PaymentMethodRepo *repositories.PaymentMethodsRepository
	MerchantRepo      *repositories.MerchantsRepository
	Gateway           services.PaymentGateway
}

func NewNicePayTransactionService(db *gorm.DB, transactionRepo *repositories.PaymentRepositoryYugabyteDB, currencyRepo *repositories.CurrenciesRepository, countryRepo *repositories.CountriesRepository, paymentMethodRepo *repositories.PaymentMethodsRepository, merchantRepo *repositories.MerchantsRepository, gateway services.PaymentGateway) *NicePayTransactionService {
	return &NicePayTransactionService{db: db, TransactionRepo: transactionRepo, CurrencyRepo: currencyRepo, CountryRepo: countryRepo, PaymentMethodRepo: paymentMethodRepo, MerchantRepo: merchantRepo, Gateway: gateway}
}

func (s *NicePayTransactionService) Save(ctx context.Context, param dto.CreatePaymentRequest, incoming entities.Incoming) (string, entities.Payment, error) {

	// find payment method
	paymentMethod, err := s.PaymentMethodRepo.FindOne(s.db, models.PaymentMethodsDataModel{Name: param.ChannelCode})
	if err != nil {
		return "", entities.Payment{}, err
	}

	currency, err := s.CurrencyRepo.FindByCode(s.db, param.Currency)
	if err != nil {
		return "", entities.Payment{}, err
	}

	country, err := s.CountryRepo.FindByCountryID(s.db, param.Country)
	if err != nil {
		return "", entities.Payment{}, err
	}
	merchant, err := s.MerchantRepo.FindOne(s.db, models.MerchantsDataModel{Name: incoming.Merchant})
	if err != nil {
		return "", entities.Payment{}, err
	}

	res, err := s.Gateway.RequestPaymentLink(ctx, nicepay.RequestPaymentLinkDTO{
		CallbackURL: configuration.AppConfig.CallbackURLNicepay,
		ReturnURL:   configuration.AppConfig.ReturnURLNicepay,
		MSISDN:      param.CustomerPhone,
		Name:        param.CustomerName,
		Number:      param.ReferenceNo,
		Channel:     param.ChannelCode,
		Amount:      param.Amount,
		Email:       param.CustomerEmail,
		Description: param.Description,
		IPAddress:   incoming.IP,
	}, configuration.AppConfig.NicepayURL)
	if err != nil {
		return "", entities.Payment{}, err
	}

	go SaveAPICall(context.Background(), &res, incoming.Merchant, err, param.ChannelCode, incoming.Path, param.CustomerPhone, incoming.Webtype, incoming.TransactionID)

	statusPending := constant.PAYMENT_STATUS_PENDING
	expiredAt := time.Now().Add(24 * time.Hour)
	err = s.TransactionRepo.Insert(s.db, &models.PaymentsDataModel{
		TransactionID:   &incoming.TransactionID,
		ReferenceNo:     &param.ReferenceNo,
		PaymentGateway:  &param.ChannelCode,
		PaymentMethodID: &paymentMethod.ID,
		CurrencyID:      &currency.ID,
		Amount:          &param.Amount,
		Description:     &param.Description,
		Status:          &statusPending,
		ExpiredPayment:  &expiredAt,
		CallbackURL:     &param.CallbackUrl,
		MerchantID:      &merchant.ID,
		CountryID:       &country.ID,
		ResponseJson:    nil,
	})
	if err != nil {
		return "", entities.Payment{}, err
	}
	// Assuming res.PaymentURL or similar exists, or just return success string?
	// Nicepay response DTO has RedirectURL
	return res.RedirectURL, entities.Payment{}, nil

}
