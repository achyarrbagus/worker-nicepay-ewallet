package repositories

import (
	"errors"
	"strings"
	"time"

	"payment-airpay/infrastructure/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MasterDataRepositoryYugabyteDB struct{}

func NewMasterDataRepositoryYugabyteDB() *MasterDataRepositoryYugabyteDB {
	return &MasterDataRepositoryYugabyteDB{}
}

func (r *MasterDataRepositoryYugabyteDB) GetOrCreateMerchant(tx *gorm.DB, code string, name string) (uuid.UUID, error) {
	if tx == nil {
		return uuid.Nil, nil
	}

	code = strings.TrimSpace(code)
	var m models.MerchantsDataModel
	err := tx.Where("code = ?", code).First(&m).Error
	if err == nil {
		return m.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
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

func (r *MasterDataRepositoryYugabyteDB) GetOrCreatePaymentMethod(tx *gorm.DB, name string) (uuid.UUID, error) {
	if tx == nil {
		return uuid.Nil, nil
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return uuid.Nil, nil
	}

	var m models.PaymentMethodsDataModel
	err := tx.Where("name = ?", name).First(&m).Error
	if err == nil {
		return m.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
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

func (r *MasterDataRepositoryYugabyteDB) GetOrCreateCurrency(tx *gorm.DB, code string) (uuid.UUID, error) {
	if tx == nil {
		return uuid.Nil, nil
	}

	code = strings.TrimSpace(code)
	if code == "" {
		return uuid.Nil, nil
	}

	var c models.CurrenciesDataModel
	err := tx.Where("code = ?", code).First(&c).Error
	if err == nil {
		return c.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
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

func (r *MasterDataRepositoryYugabyteDB) GetOrCreateVAProvider(tx *gorm.DB, name string, providerName string) (uuid.UUID, error) {
	if tx == nil {
		return uuid.Nil, nil
	}

	providerName = strings.TrimSpace(providerName)
	if providerName == "" {
		return uuid.Nil, nil
	}

	var p models.VAProvidersDataModel
	err := tx.Where("provider_name = ?", providerName).First(&p).Error
	if err == nil {
		return p.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = providerName
	}

	actor := "system"
	dataStatus := "ACTIVE"
	now := time.Now().UnixMilli()
	newP := models.VAProvidersDataModel{
		Name:         name,
		ProviderName: providerName,
		CreatedDate:  &now,
		CreatedUser:  &actor,
		DataStatus:   &dataStatus,
	}
	if err := tx.Create(&newP).Error; err != nil {
		return uuid.Nil, err
	}
	return newP.ID, nil
}

func (r *MasterDataRepositoryYugabyteDB) GetOrCreateEWalletProvider(tx *gorm.DB, providerName string) (uuid.UUID, error) {
	if tx == nil {
		return uuid.Nil, nil
	}

	providerName = strings.TrimSpace(providerName)
	if providerName == "" {
		return uuid.Nil, nil
	}

	var p models.EWalletProvidersDataModel
	err := tx.Where("provider_name = ?", providerName).First(&p).Error
	if err == nil {
		return p.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, err
	}

	actor := "system"
	dataStatus := "ACTIVE"
	now := time.Now().UnixMilli()
	newP := models.EWalletProvidersDataModel{
		Name:         providerName,
		ProviderName: providerName,
		CreatedDate:  &now,
		CreatedUser:  &actor,
		DataStatus:   &dataStatus,
	}
	if err := tx.Create(&newP).Error; err != nil {
		return uuid.Nil, err
	}
	return newP.ID, nil
}

func (r *MasterDataRepositoryYugabyteDB) GetOrCreateCountry(tx *gorm.DB, countryID string) (uuid.UUID, error) {
	if tx == nil {
		return uuid.Nil, nil
	}

	countryID = strings.TrimSpace(countryID)
	if countryID == "" {
		return uuid.Nil, nil
	}

	var c models.CountriesDataModel
	err := tx.Where("country_id = ?", countryID).First(&c).Error
	if err == nil {
		return c.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
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
