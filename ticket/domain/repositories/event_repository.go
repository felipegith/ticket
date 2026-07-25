package repositories

import (
	"context"
	"ticket/ticket/domain/entities"
)

type EventRepository interface {
	Create(ctx context.Context, event *entities.Event) error
	GetAll(ctx context.Context) ([]*entities.Event, error)
	GetByID(ctx context.Context, id string) (*entities.Event, error)
}
