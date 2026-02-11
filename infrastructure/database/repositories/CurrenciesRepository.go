package repositories

import (
	"errors"
	"worker-nicepay/infrastructure/database/models"

	"gorm.io/gorm"
)

type CurrenciesRepository struct{}

func NewCurrenciesRepository() *CurrenciesRepository {
	return &CurrenciesRepository{}
}

func (r *CurrenciesRepository) Insert(tx *gorm.DB, model *models.CurrenciesDataModel) error {
	if tx == nil || model == nil {
		return nil
	}
	return tx.Create(model).Error
}

func (r *CurrenciesRepository) FindAll(tx *gorm.DB) ([]models.CurrenciesDataModel, error) {
	if tx == nil {
		return nil, nil
	}
	var currencies []models.CurrenciesDataModel
	err := tx.Find(&currencies).Error
	return currencies, err
}

func (r *CurrenciesRepository) FindByCode(tx *gorm.DB, code string) (*models.CurrenciesDataModel, error) {
	if tx == nil {
		return nil, nil
	}
	var currency models.CurrenciesDataModel
	err := tx.Where("code = ?", code).First(&currency).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &currency, nil
}
