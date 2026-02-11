package repositories

import (
	"errors"
	"worker-nicepay/infrastructure/database/models"

	"gorm.io/gorm"
)

type CountriesRepository struct{}

func NewCountriesRepository() *CountriesRepository {
	return &CountriesRepository{}
}

func (r *CountriesRepository) Insert(tx *gorm.DB, model *models.CountriesDataModel) error {
	if tx == nil || model == nil {
		return nil
	}
	return tx.Create(model).Error
}

func (r *CountriesRepository) FindAll(tx *gorm.DB) ([]models.CountriesDataModel, error) {
	if tx == nil {
		return nil, nil
	}
	var countries []models.CountriesDataModel
	err := tx.Find(&countries).Error
	return countries, err
}

func (r *CountriesRepository) FindByCountryID(tx *gorm.DB, countryID string) (*models.CountriesDataModel, error) {
	if tx == nil {
		return nil, nil
	}
	var country models.CountriesDataModel
	err := tx.Where("country_id = ?", countryID).First(&country).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &country, nil
}
