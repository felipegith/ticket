package repositories

import (
	"context"
	"ticket/ticket/domain/entities"
)

type OutboxRepository interface {
	FetchPending(ctx context.Context, limit int) ([]*entities.Outbox, error)

	MarkAsProcessed(ctx context.Context, message *entities.Outbox) error
}
