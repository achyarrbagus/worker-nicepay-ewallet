package repositories

import (
	"payment-airpay/infrastructure/database/models"

	"gorm.io/gorm"
)

type PaymentRepositoryYugabyteDB struct{}

func NewPaymentRepositoryYugabyteDB() *PaymentRepositoryYugabyteDB {
	return &PaymentRepositoryYugabyteDB{}
}

func (r *PaymentRepositoryYugabyteDB) Insert(tx *gorm.DB, model *models.PaymentsDataModel) error {
	if tx == nil || model == nil {
		return nil
	}
	return tx.Create(model).Error
}
