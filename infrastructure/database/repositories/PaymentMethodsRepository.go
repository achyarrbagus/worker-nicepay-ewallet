package repositories

import (
	"worker-nicepay/infrastructure/database/models"

	"gorm.io/gorm"
)

type PaymentMethodsRepository struct{}

func NewPaymentMethodsRepository() *PaymentMethodsRepository {
	return &PaymentMethodsRepository{}
}

func (r *PaymentMethodsRepository) Insert(tx *gorm.DB, model *models.PaymentMethodsDataModel) error {
	if tx == nil || model == nil {
		return nil
	}
	return tx.Create(model).Error
}

func (r *PaymentMethodsRepository) FindAll(tx *gorm.DB) ([]models.PaymentMethodsDataModel, error) {
	if tx == nil {
		return nil, nil
	}
	var methods []models.PaymentMethodsDataModel
	err := tx.Find(&methods).Error
	return methods, err
}

func (r *PaymentMethodsRepository) FindOne(tx *gorm.DB, where models.PaymentMethodsDataModel) (models.PaymentMethodsDataModel, error) {
	if tx == nil {
		return models.PaymentMethodsDataModel{}, nil
	}
	var method models.PaymentMethodsDataModel
	err := tx.Where(&where).Last(&method).Error
	return method, err
}
