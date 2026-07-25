package repositories

import (
	"context"
	"ticket/ticket/domain/entities"
)

type TicketRepository interface {
	Create(ctx context.Context, ticket *entities.Ticket) error
	SeatExists(ctx context.Context, eventID, seat string) (bool, error)
	Update(ctx context.Context, ticketId string) (bool, error)
}
