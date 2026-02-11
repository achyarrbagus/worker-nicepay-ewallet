package repositories

import (
	"errors"
	"worker-nicepay/infrastructure/database/models"

	"gorm.io/gorm"
)

type MerchantsRepository struct {
}

func NewMerchantsRepository() *MerchantsRepository {
	return &MerchantsRepository{}
}

func (r *MerchantsRepository) Insert(tx *gorm.DB, model *models.MerchantsDataModel) error {
	if tx == nil || model == nil {
		return nil
	}
	return tx.Create(model).Error
}

func (r *MerchantsRepository) FindAll(tx *gorm.DB) ([]models.MerchantsDataModel, error) {
	if tx == nil {
		return nil, nil
	}
	var merchants []models.MerchantsDataModel
	err := tx.Find(&merchants).Error
	return merchants, err
}

func (r *MerchantsRepository) FindOne(tx *gorm.DB, where models.MerchantsDataModel) (*models.MerchantsDataModel, error) {
	if tx == nil {
		return nil, nil
	}
	var merchant models.MerchantsDataModel
	err := tx.Where(&where).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &merchant, nil
}
