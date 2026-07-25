package repositories

import "context"

type BookingRepository interface {
	Create(ctx context.Context, ticketId, seat, userId string) error
}
