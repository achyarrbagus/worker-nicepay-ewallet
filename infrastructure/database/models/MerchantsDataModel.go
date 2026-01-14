package models

import "github.com/google/uuid"

type MerchantsDataModel struct {
	ID          uuid.UUID `gorm:"primaryKey;column:id;type:uuid"`
	Name        string    `gorm:"column:name"`
	Code        string    `gorm:"column:code;uniqueIndex"`
	Username    string    `gorm:"column:username"`
	Password    string    `gorm:"column:password"`
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
