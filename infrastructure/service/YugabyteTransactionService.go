package service

import (
	"context"
	"encoding/json"
	"time"

	"payment-airpay/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		_ = db.AutoMigrate(&YugabyteTransaction{})
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
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	record := YugabyteTransaction{
		TransactionID: payment.PaymentRequestID,
		Payload:       string(payloadBytes),
		CreatedAt:     time.Now(),
	}

	// If transaction_id already exists (mock often returns same ID), update the payload/created_at
	return s.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "transaction_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"payload", "created_at"}),
		}).
		Create(&record).Error
}
