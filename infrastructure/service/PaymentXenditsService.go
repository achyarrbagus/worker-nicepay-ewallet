package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"payment-airpay/domain/entities"
	"payment-airpay/infrastructure/database/clients"
	"payment-airpay/infrastructure/database/models"
	"payment-airpay/infrastructure/database/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentXendit struct {
	masterRepo  *repositories.MasterDataRepositoryYugabyteDB
	paymentRepo *repositories.PaymentRepositoryYugabyteDB
	xenditRepo  *repositories.XenditRepositoryYugabyteDB
	db          clients.YugabyteClient
}

func NewPaymentXendit(
	masterRepo *repositories.MasterDataRepositoryYugabyteDB,
	paymentRepo *repositories.PaymentRepositoryYugabyteDB,
	xenditRepo *repositories.XenditRepositoryYugabyteDB,
	db clients.YugabyteClient,
) *PaymentXendit {
	return &PaymentXendit{masterRepo: masterRepo, paymentRepo: paymentRepo, xenditRepo: xenditRepo, db: db}
}

func (s *PaymentXendit) Save(ctx context.Context, payment entities.Payment, payload map[string]interface{}) error {
	if s == nil || s.db == nil || s.db.GetDB() == nil {
		return nil
	}

	db := s.db.GetDB().WithContext(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		merchantName, paymentGateway := extractMetadata(payload)
		merchantCode := extractMerchantCode(payload, merchantName)
		channelCode := payment.ChannelCode
		currencyCode := payment.Currency
		countryCode := payment.Country

		merchantID, err := s.masterRepo.GetOrCreateMerchant(tx, merchantCode, merchantName)
		if err != nil {
			return err
		}
		paymentMethodName := normalizePaymentMethodName(channelCode)
		paymentMethodID, err := s.masterRepo.GetOrCreatePaymentMethod(tx, paymentMethodName)
		if err != nil {
			return err
		}
		currencyID, err := s.masterRepo.GetOrCreateCurrency(tx, currencyCode)
		if err != nil {
			return err
		}
		countryID, err := s.masterRepo.GetOrCreateCountry(tx, countryCode)
		if err != nil {
			return err
		}

		var merchantIDPtr *uuid.UUID
		if merchantID != uuid.Nil {
			merchantIDPtr = &merchantID
		}
		var paymentMethodIDPtr *uuid.UUID
		if paymentMethodID != uuid.Nil {
			paymentMethodIDPtr = &paymentMethodID
		}
		var currencyIDPtr *uuid.UUID
		if currencyID != uuid.Nil {
			currencyIDPtr = &currencyID
		}
		var countryIDPtr *uuid.UUID
		if countryID != uuid.Nil {
			countryIDPtr = &countryID
		}

		status := string(payment.Status)
		amount := payment.RequestAmount
		desc := payment.Description
		trxID := payment.PaymentRequestID
		refNo := payment.ReferenceID
		dataStatus := "ACTIVE"
		actor := "system"

		var createdDate *int64
		if strings.TrimSpace(payment.Created) != "" {
			if t, err := time.Parse(time.RFC3339Nano, payment.Created); err == nil {
				v := t.UnixMilli()
				createdDate = &v
			}
		}
		var updatedDate *int64
		if strings.TrimSpace(payment.Updated) != "" {
			if t, err := time.Parse(time.RFC3339Nano, payment.Updated); err == nil {
				v := t.UnixMilli()
				updatedDate = &v
			}
		}

		var expiredAt *time.Time
		if v, ok := payment.ChannelProps["expires_at"]; ok {
			if s, ok := v.(string); ok {
				if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
					expiredAt = &t
				}
			}
		}

		responseJSON, err := json.Marshal(payment)
		if err != nil {
			return err
		}

		// Persist Xendit gateway-specific payload
		if isEWalletPayment(payment, payload) {
			providerName := strings.ToUpper(strings.TrimSpace(channelCode))
			providerID, err := s.masterRepo.GetOrCreateEWalletProvider(tx, providerName)
			if err != nil {
				return err
			}

			var providerIDPtr *uuid.UUID
			if providerID != uuid.Nil {
				providerIDPtr = &providerID
			}

			var customerMSISDN *string
			if payment.ChannelProps != nil {
				if v, ok := payment.ChannelProps["account_mobile_number"]; ok {
					if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
						t := s
						customerMSISDN = &t
					}
				}
			}

			px := models.PaymentXenditEWalletsDataModel{
				TransactionID:     &trxID,
				URLReturn:         "",
				EWalletProviderID: providerIDPtr,
				CustomerMSISDN:    customerMSISDN,
				ResponseJson:      responseJSON,
				CreatedDate:       createdDate,
				CreatedUser:       &actor,
				UpdatedDate:       updatedDate,
				UpdatedUser:       &actor,
				DataStatus:        &dataStatus,
			}
			if err := s.xenditRepo.InsertEWallets(tx, &px); err != nil {
				return err
			}
		} else if strings.Contains(strings.ToUpper(channelCode), "VIRTUAL_ACCOUNT") {
			providerName := extractVAProviderName(channelCode)
			providerID, err := s.masterRepo.GetOrCreateVAProvider(tx, channelCode, providerName)
			if err != nil {
				return err
			}

			var providerIDPtr *uuid.UUID
			if providerID != uuid.Nil {
				providerIDPtr = &providerID
			}

			vaNumber := extractVANumber(payment)
			px := models.PaymentXenditVasDataModel{
				TransactionID: &trxID,
				URLReturn:     "",
				VAProviderID:  providerIDPtr,
				VANumber:      vaNumber,
				ResponseJson:  responseJSON,
				CreatedDate:   createdDate,
				CreatedUser:   &actor,
				UpdatedDate:   updatedDate,
				UpdatedUser:   &actor,
				DataStatus:    &dataStatus,
			}
			if err := s.xenditRepo.InsertVAs(tx, &px); err != nil {
				return err
			}
		} else if strings.EqualFold(channelCode, "QRIS") {
			qrpyID := extractQRPYID(payment)
			px := models.PaymentXenditQrisesDataModel{
				TransactionID: &trxID,
				URLReturn:     "",
				QRPYID:        qrpyID,
				ResponseJson:  responseJSON,
				CreatedDate:   createdDate,
				CreatedUser:   &actor,
				UpdatedDate:   updatedDate,
				UpdatedUser:   &actor,
				DataStatus:    &dataStatus,
			}
			if err := s.xenditRepo.InsertQrises(tx, &px); err != nil {
				return err
			}
		}

		model := models.PaymentsDataModel{
			TransactionID:   &trxID,
			PaymentGateway:  &paymentGateway,
			ReferenceNo:     &refNo,
			PaymentMethodID: paymentMethodIDPtr,
			CurrencyID:      currencyIDPtr,
			Amount:          &amount,
			Description:     &desc,
			Status:          &status,
			ExpiredPayment:  expiredAt,
			MerchantID:      merchantIDPtr,
			CountryID:       countryIDPtr,
			ResponseJson:    responseJSON,
			CreatedDate:     createdDate,
			CreatedUser:     &actor,
			UpdatedDate:     updatedDate,
			UpdatedUser:     &actor,
			DataStatus:      &dataStatus,
		}

		return s.paymentRepo.Insert(tx, &model)
	})
}

func normalizePaymentMethodName(channelCode string) string {
	cc := strings.ToUpper(strings.TrimSpace(channelCode))
	if strings.Contains(cc, "VIRTUAL_ACCOUNT") {
		return "VA"
	}
	if strings.EqualFold(cc, "QRIS") {
		return "QRIS"
	}
	return "E_WALLET"
}

func isEWalletPayment(payment entities.Payment, payload map[string]interface{}) bool {
	if channelProps, ok := payload["channel_properties"].(map[string]interface{}); ok {
		if _, has := channelProps["account_mobile_number"]; has {
			return true
		}
	}
	if payment.ChannelProps != nil {
		if _, has := payment.ChannelProps["account_mobile_number"]; has {
			return true
		}
	}
	return false
}

func extractMetadata(payload map[string]interface{}) (merchantName string, paymentGateway string) {
	paymentGateway = ""
	if payload == nil {
		return "", paymentGateway
	}
	meta, _ := payload["metadata"].(map[string]interface{})
	if meta == nil {
		return "", paymentGateway
	}
	if v, ok := meta["merchant"].(string); ok {
		merchantName = v
	}
	if v, ok := meta["payment_gateway"].(string); ok {
		paymentGateway = v
	}
	return merchantName, paymentGateway
}

func extractQRPYID(payment entities.Payment) string {
	for _, a := range payment.Actions {
		if strings.EqualFold(a.Descriptor, "QR_STRING") && strings.TrimSpace(a.Value) != "" {
			return a.Value
		}
	}
	for _, a := range payment.Actions {
		if strings.EqualFold(a.Type, "PRESENT_TO_CUSTOMER") && strings.TrimSpace(a.Value) != "" {
			return a.Value
		}
	}
	return ""
}

func extractVANumber(payment entities.Payment) string {
	for _, a := range payment.Actions {
		if strings.EqualFold(a.Descriptor, "VIRTUAL_ACCOUNT_NUMBER") && strings.TrimSpace(a.Value) != "" {
			return a.Value
		}
	}
	for _, a := range payment.Actions {
		if strings.EqualFold(a.Type, "PRESENT_TO_CUSTOMER") && strings.TrimSpace(a.Value) != "" {
			return a.Value
		}
	}
	return ""
}

func extractVAProviderName(channelCode string) string {
	cc := strings.ToUpper(strings.TrimSpace(channelCode))
	cc = strings.TrimSuffix(cc, "_VIRTUAL_ACCOUNT")
	return cc
}

func extractMerchantCode(payload map[string]interface{}, merchantName string) string {
	// Prefer explicit code if provided
	if payload != nil {
		if v, ok := payload["merchant_code"].(string); ok && v != "" {
			return v
		}
	}
	// Fallback: first token of name ("PPOB VIA" -> "PPOB")
	merchantName = strings.TrimSpace(merchantName)
	if merchantName == "" {
		return "PPOB"
	}
	parts := strings.Fields(merchantName)
	if len(parts) > 0 {
		return parts[0]
	}
	return "PPOB"
}

func getOrCreateMerchant(tx *gorm.DB, code string, name string) (uuid.UUID, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		code = ""
	}

	var m models.MerchantsDataModel
	err := tx.Where("code = ?", code).First(&m).Error
	if err == nil {
		return m.ID, nil
	}
	if !errorsIsRecordNotFound(err) {
		return uuid.Nil, err
	}

	if strings.TrimSpace(name) == "" {
		name = code
	}

	actor := "system"
	dataStatus := "ACTIVE"
	now := time.Now().UnixMilli()
	newM := models.MerchantsDataModel{
		Code:        code,
		Name:        name,
		CreatedDate: &now,
		CreatedUser: &actor,
		DataStatus:  &dataStatus,
	}
	if err := tx.Create(&newM).Error; err != nil {
		return uuid.Nil, err
	}
	return newM.ID, nil
}

func getOrCreatePaymentMethod(tx *gorm.DB, name string) (uuid.UUID, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return uuid.Nil, nil
	}

	var m models.PaymentMethodsDataModel
	err := tx.Where("name = ?", name).First(&m).Error
	if err == nil {
		return m.ID, nil
	}
	if !errorsIsRecordNotFound(err) {
		return uuid.Nil, err
	}

	actor := "system"
	dataStatus := "ACTIVE"
	now := time.Now().UnixMilli()
	newM := models.PaymentMethodsDataModel{
		Name:        name,
		CreatedDate: &now,
		CreatedUser: &actor,
		DataStatus:  &dataStatus,
	}
	if err := tx.Create(&newM).Error; err != nil {
		return uuid.Nil, err
	}
	return newM.ID, nil
}

func getOrCreateCurrency(tx *gorm.DB, code string) (uuid.UUID, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return uuid.Nil, nil
	}

	var c models.CurrenciesDataModel
	err := tx.Where("code = ?", code).First(&c).Error
	if err == nil {
		return c.ID, nil
	}
	if !errorsIsRecordNotFound(err) {
		return uuid.Nil, err
	}

	name := code
	if strings.EqualFold(code, "IDR") {
		name = "Indonesia"
	}
	actor := "system"
	dataStatus := "ACTIVE"
	now := time.Now().UnixMilli()
	newC := models.CurrenciesDataModel{
		Code:        code,
		Name:        name,
		CreatedDate: &now,
		CreatedUser: &actor,
		DataStatus:  &dataStatus,
	}
	if err := tx.Create(&newC).Error; err != nil {
		return uuid.Nil, err
	}
	return newC.ID, nil
}

func getOrCreateCountry(tx *gorm.DB, countryID string) (uuid.UUID, error) {
	countryID = strings.TrimSpace(countryID)
	if countryID == "" {
		return uuid.Nil, nil
	}

	var c models.CountriesDataModel
	err := tx.Where("country_id = ?", countryID).First(&c).Error
	if err == nil {
		return c.ID, nil
	}
	if !errorsIsRecordNotFound(err) {
		return uuid.Nil, err
	}

	name := countryID
	if strings.EqualFold(countryID, "ID") {
		name = "Indonesia"
	}
	actor := "system"
	dataStatus := "ACTIVE"
	now := time.Now().UnixMilli()
	newC := models.CountriesDataModel{
		CountryId:   countryID,
		Name:        name,
		CreatedDate: &now,
		CreatedUser: &actor,
		DataStatus:  &dataStatus,
	}
	if err := tx.Create(&newC).Error; err != nil {
		return uuid.Nil, err
	}
	return newC.ID, nil
}

func errorsIsRecordNotFound(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrRecordNotFound)
}
