package models

import "github.com/google/uuid"

type PaymentMethodsDataModel struct {
	ID          uuid.UUID `gorm:"primaryKey;column:id;type:uuid"`
	Name        string    `gorm:"column:name;uniqueIndex"`
	CreatedDate *int64
	CreatedUser *string
	CreatedIp   *string
	UpdatedDate *int64
	UpdatedUser *string
	UpdatedIp   *string
	DeletedDate *int64
	DeletedUser *string
	DeletedIp   *string
	DataStatus  *string
}
