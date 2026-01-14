package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PaymentsDataModel struct {
	ID              uuid.UUID                `gorm:"primaryKey;column:id;type:uuid"`
	TransactionID   *string                  `gorm:"column:transaction_id"`
	PaymentGateway  *string                  `gorm:"column:payment_gateway"`
	ReferenceNo     *string                  `gorm:"column:reference_no"`
	PaymentMethodID *uuid.UUID               `gorm:"column:payment_method_id;type:uuid"`
	CurrencyID      *uuid.UUID               `gorm:"column:currency_id;type:uuid"`
	Amount          *float64                 `gorm:"column:amount"`
	Description     *string                  `gorm:"column:description"`
	Status          *string                  `gorm:"column:status"`
	ExpiredPayment  *time.Time               `gorm:"column:expired_payment"`
	CallbackURL     *string                  `gorm:"column:callback_url"`
	MerchantID      *uuid.UUID               `gorm:"column:merchant_id;type:uuid"`
	CountryID       *uuid.UUID               `gorm:"column:country_id;type:uuid"`
	Merchant        *MerchantsDataModel      `gorm:"foreignKey:MerchantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	PaymentMethod   *PaymentMethodsDataModel `gorm:"foreignKey:PaymentMethodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Currency        *CurrenciesDataModel     `gorm:"foreignKey:CurrencyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Country         *CountriesDataModel      `gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ResponseJson    json.RawMessage          `gorm:"column:response_json;type:jsonb"`
	CreatedDate     *int64
	CreatedUser     *string
	CreatedIp       *string
	UpdatedDate     *int64
	UpdatedUser     *string
	UpdatedIp       *string
	DeletedDate     *int64
	DeletedUser     *string
	DeletedIp       *string
	DataStatus      *string
}
