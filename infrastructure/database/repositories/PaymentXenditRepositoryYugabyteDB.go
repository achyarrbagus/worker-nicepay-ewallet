package repositories

import (
	"context"

	"payment-airpay/infrastructure/database/clients"
	"payment-airpay/infrastructure/database/models"

	"gorm.io/gorm/clause"
)

type PaymentXenditRepositoryYugabyteDB struct {
	db clients.YugabyteClient
}

func NewPaymentXenditRepositoryYugabyteDB(db clients.YugabyteClient) *PaymentXenditRepositoryYugabyteDB {
	return &PaymentXenditRepositoryYugabyteDB{db: db}
}

func (r *PaymentXenditRepositoryYugabyteDB) Upsert(ctx context.Context) error {
	if r == nil || r.db == nil || r.db.GetDB() == nil {
		return nil
	}

	model := models.PaymentXenditQrisesDataModel{
		TransactionID: nil,
		URLReturn:     "",
		QRPYID:        "",
	}

	// Placeholder upsert: unique by transaction_id; if it's NULL, each insert becomes a new row.
	// This matches the request: keep empty values for now.
	return r.db.GetDB().WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&model).Error
}
