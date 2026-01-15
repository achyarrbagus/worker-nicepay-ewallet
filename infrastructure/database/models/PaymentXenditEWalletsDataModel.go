package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type PaymentXenditEWalletsDataModel struct {
	ID                uuid.UUID                  `gorm:"primaryKey;column:id;type:uuid"`
	TransactionID     *string                    `gorm:"column:transaction_id;uniqueIndex"`
	URLReturn         string                     `gorm:"column:url_return"`
	EWalletProviderID *uuid.UUID                 `gorm:"column:e_wallet_provider_id;type:uuid"`
	EWalletProvider   *EWalletProvidersDataModel `gorm:"foreignKey:EWalletProviderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CustomerUsername  *string                    `gorm:"column:customer_username"`
	CustomerMSISDN    *string                    `gorm:"column:customer_msisdn"`
	CustomerEmail     *string                    `gorm:"column:customer_email"`
	ResponseURL       *string                    `gorm:"column:response_url"`
	ResponseJson      json.RawMessage            `gorm:"column:response_json;type:jsonb"`
	CreatedDate       *int64
	CreatedUser       *string
	CreatedIp         *string
	UpdatedDate       *int64
	UpdatedUser       *string
	UpdatedIp         *string
	DeletedDate       *int64
	DeletedUser       *string
	DeletedIp         *string
	DataStatus        *string
}

func (PaymentXenditEWalletsDataModel) TableName() string {
	return "payment_xendit_ewallets"
}
