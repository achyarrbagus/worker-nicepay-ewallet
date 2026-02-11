package repositories

import (
	"worker-nicepay/infrastructure/database/models"

	"gorm.io/gorm"
)

type EWalletProvidersRepository struct{}

func NewEWalletProvidersRepository() *EWalletProvidersRepository {
	return &EWalletProvidersRepository{}
}

func (r *EWalletProvidersRepository) Insert(tx *gorm.DB, model *models.EWalletProvidersDataModel) error {
	if tx == nil || model == nil {
		return nil
	}
	return tx.Create(model).Error
}

func (r *EWalletProvidersRepository) FindAll(tx *gorm.DB) ([]models.EWalletProvidersDataModel, error) {
	if tx == nil {
		return nil, nil
	}
	var providers []models.EWalletProvidersDataModel
	err := tx.Find(&providers).Error
	return providers, err
}
