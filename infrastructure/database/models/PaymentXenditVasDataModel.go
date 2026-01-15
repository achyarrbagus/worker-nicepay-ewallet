package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type PaymentXenditVasDataModel struct {
	ID               uuid.UUID             `gorm:"primaryKey;column:id;type:uuid"`
	TransactionID    *string               `gorm:"column:transaction_id;uniqueIndex"`
	URLReturn        string                `gorm:"column:url_return"`
	VAProviderID     *uuid.UUID            `gorm:"column:va_provider_id;type:uuid"`
	VAProvider       *VAProvidersDataModel `gorm:"foreignKey:VAProviderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	VANumber         string                `gorm:"column:va_number"`
	CustomerUsername *string               `gorm:"column:customer_username"`
	CustomerMSISDN   *string               `gorm:"column:customer_msisdn"`
	CustomerEmail    *string               `gorm:"column:customer_email"`
	ResponseURL      *string               `gorm:"column:response_url"`
	ResponseJson     json.RawMessage       `gorm:"column:response_json;type:jsonb"`
	CreatedDate      *int64
	CreatedUser      *string
	CreatedIp        *string
	UpdatedDate      *int64
	UpdatedUser      *string
	UpdatedIp        *string
	DeletedDate      *int64
	DeletedUser      *string
	DeletedIp        *string
	DataStatus       *string
}
