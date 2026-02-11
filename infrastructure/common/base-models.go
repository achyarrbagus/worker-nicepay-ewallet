package common

import (
	"time"

	"gorm.io/gorm"
)

type BaseModelSoftDelete struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
