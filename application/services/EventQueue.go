package services

import "context"

type EventQueue interface {
	Enqueue(ctx context.Context, event Event) error
}
