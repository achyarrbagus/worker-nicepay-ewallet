package service

import (
	"context"
	"log"
)

type LogTransactionService struct{}

func NewLogTransactionService() *LogTransactionService { return &LogTransactionService{} }

func (s *LogTransactionService) Save(ctx context.Context, transactionID string, payload map[string]interface{}) error {
	if transactionID == "" {
		log.Println("transaction save skipped: empty transaction_id")
		return nil
	}
	log.Printf("transaction saved: id=%s payload=%v\n", transactionID, payload)
	return nil
}
