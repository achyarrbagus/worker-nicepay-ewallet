package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"payment-airpay/domain/entities"
	"payment-airpay/infrastructure/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type YugabyteTransaction struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid"`
	TransactionID string    `gorm:"column:transaction_id;uniqueIndex"`
	Payload       string    `gorm:"column:payload;type:text"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (m *YugabyteTransaction) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		m.ID = id
	}
	return nil
}

type YugabyteTransactionService struct {
	db *gorm.DB
}

func NewYugabyteTransactionService(db *gorm.DB) *YugabyteTransactionService {
	if db != nil {
		_ = db.AutoMigrate(&YugabyteTransaction{}, &models.PaymentXenditEWalletsDataModel{}, &models.PaymentXenditQrisesDataModel{}, &models.PaymentXenditVasDataModel{}, &models.EWalletProvidersDataModel{})
	}
	return &YugabyteTransactionService{db: db}
}

func (s *YugabyteTransactionService) Save(ctx context.Context, payment entities.Payment, payload map[string]interface{}) error {
	if s == nil || s.db == nil {
		return nil
	}
	if payment.PaymentRequestID == "" {
		return nil
	}

	// Check if this is an e-wallet payment by looking for account_mobile_number in channel_properties
	if isEWalletPaymentTxSvc(payment, payload) {
		return s.saveToEWalletTable(ctx, payment)
	} else if isQRISPayment(payment) {
		return s.saveToQRISTable(ctx, payment, payload)
	} else {
		return s.saveToVATable(ctx, payment, payload)
	}
}

func isEWalletPaymentTxSvc(payment entities.Payment, payload map[string]interface{}) bool {
	// Check if channel_properties contains account_mobile_number
	if channelProps, ok := payload["channel_properties"].(map[string]interface{}); ok {
		if _, hasMobileNumber := channelProps["account_mobile_number"]; hasMobileNumber {
			return true
		}
	}
	// Also check response channel properties
	if payment.ChannelProps != nil {
		if _, hasMobileNumber := payment.ChannelProps["account_mobile_number"]; hasMobileNumber {
			return true
		}
	}
	return false
}

func isQRISPayment(payment entities.Payment) bool {
	return strings.ToUpper(payment.ChannelCode) == "QRIS"
}

func (s *YugabyteTransactionService) saveToEWalletTable(ctx context.Context, payment entities.Payment) error {
	responseJSON, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	// Find ewallet provider by channel code
	var ewalletProvider models.EWalletProvidersDataModel
	if err := s.db.Where("provider_name = ?", payment.ChannelCode).First(&ewalletProvider).Error; err != nil {
		// Create new provider if not found
		ewalletProvider = models.EWalletProvidersDataModel{
			Name:         payment.ChannelCode,
			ProviderName: payment.ChannelCode,
		}
		if err := s.db.Create(&ewalletProvider).Error; err != nil {
			return err
		}
	}

	record := models.PaymentXenditEWalletsDataModel{
		TransactionID:     &payment.PaymentRequestID,
		ResponseJson:      responseJSON,
		EWalletProviderID: &ewalletProvider.ID,
	}

	// Set customer mobile number if available
	if payment.ChannelProps != nil {
		if mobileNumber, ok := payment.ChannelProps["account_mobile_number"].(string); ok && mobileNumber != "" {
			record.CustomerMSISDN = &mobileNumber
		}
	}

	return s.db.WithContext(ctx).Create(&record).Error
}

func (s *YugabyteTransactionService) saveToQRISTable(ctx context.Context, payment entities.Payment, payload map[string]interface{}) error {
	responseJSON, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	record := models.PaymentXenditQrisesDataModel{
		TransactionID: &payment.PaymentRequestID,
		ResponseJson:  responseJSON,
	}

	return s.db.WithContext(ctx).Create(&record).Error
}

func (s *YugabyteTransactionService) saveToVATable(ctx context.Context, payment entities.Payment, payload map[string]interface{}) error {
	responseJSON, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	record := models.PaymentXenditVasDataModel{
		TransactionID: &payment.PaymentRequestID,
		ResponseJson:  responseJSON,
	}

	return s.db.WithContext(ctx).Create(&record).Error
}
