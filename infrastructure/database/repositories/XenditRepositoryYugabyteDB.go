package repositories

import (
	"payment-airpay/infrastructure/database/models"

	"gorm.io/gorm"
)

type XenditRepositoryYugabyteDB struct{}

func NewXenditRepositoryYugabyteDB() *XenditRepositoryYugabyteDB {
	return &XenditRepositoryYugabyteDB{}
}

func (r *XenditRepositoryYugabyteDB) InsertQrises(tx *gorm.DB, model *models.PaymentXenditQrisesDataModel) error {
	if tx == nil || model == nil {
		return nil
	}
	return tx.Create(model).Error
}

func (r *XenditRepositoryYugabyteDB) InsertVAs(tx *gorm.DB, model *models.PaymentXenditVasDataModel) error {
	if tx == nil || model == nil {
		return nil
	}
	return tx.Create(model).Error
}

func (r *XenditRepositoryYugabyteDB) InsertEWallets(tx *gorm.DB, model *models.PaymentXenditEWalletsDataModel) error {
	if tx == nil || model == nil {
		return nil
	}
	return tx.Create(model).Error
}
