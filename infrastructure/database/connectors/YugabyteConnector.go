package connectors

import (
	"context"

	"gorm.io/gorm"
)

type YugabyteConnector struct {
	client *gorm.DB
}

func NewYugabyteConnector(client *gorm.DB) *YugabyteConnector {
	return &YugabyteConnector{client: client}
}

func (c *YugabyteConnector) Create(ctx context.Context, value interface{}) error {
	return c.client.WithContext(ctx).Create(value).Error
}

func (c *YugabyteConnector) First(ctx context.Context, dest interface{}, query interface{}, args ...interface{}) error {
	if query == nil {
		return c.client.WithContext(ctx).First(dest, args...).Error
	}
	conds := append([]interface{}{query}, args...)
	return c.client.WithContext(ctx).First(dest, conds...).Error
}

func (c *YugabyteConnector) Save(ctx context.Context, value interface{}) error {
	return c.client.WithContext(ctx).Save(value).Error
}

func (c *YugabyteConnector) GetDB() *gorm.DB {
	return c.client
}
