package models

import "github.com/google/uuid"

type EWalletProvidersDataModel struct {
	ID           uuid.UUID `gorm:"primaryKey;column:id;type:uuid"`
	Name         string    `gorm:"column:name"`
	ProviderName string    `gorm:"column:provider_name;uniqueIndex"`
	CreatedDate  *int64
	CreatedUser  *string
	CreatedIp    *string
	UpdatedDate  *int64
	UpdatedUser  *string
	UpdatedIp    *string
	DeletedDate  *int64
	DeletedUser  *string
	DeletedIp    *string
	DataStatus   *string
}

func (EWalletProvidersDataModel) TableName() string {
	return "ewallet_providers"
}
