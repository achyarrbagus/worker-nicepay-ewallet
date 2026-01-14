package clients

import (
	"context"

	"gorm.io/gorm"
)

type YugabyteClient interface {
	Create(ctx context.Context, value interface{}) error
	First(ctx context.Context, dest interface{}, query interface{}, args ...interface{}) error
	Save(ctx context.Context, value interface{}) error
	GetDB() *gorm.DB
}
