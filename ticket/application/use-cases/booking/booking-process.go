package booking

import (
	"context"
	"ticket/ticket/application/ports"
	"ticket/ticket/domain/repositories"
)

type BookingProcess struct {
	ticketRepository  repositories.TicketRepository
	bookingRepository repositories.BookingRepository
	redis             ports.Cache
}

func NewBookingProcess(ticketRepository repositories.TicketRepository, bookingRepository repositories.BookingRepository, redis ports.Cache) *BookingProcess {
	return &BookingProcess{
		ticketRepository:  ticketRepository,
		bookingRepository: bookingRepository,
		redis:             redis,
	}
}

func (b *BookingProcess) Process(ctx context.Context, message *BookingInputModel) error {
	return b.Execute(ctx, message.TicketId, message.Seat, message.UserId)
}

func (b *BookingProcess) Execute(ctx context.Context, ticketId, seat, userId string) error {

	hasUpate, err := b.ticketRepository.Update(ctx, ticketId)

	if err != nil {
		return err
	}

	if hasUpate {
		err := b.bookingRepository.Create(ctx, ticketId, seat, userId)
		if err != nil {
			return err
		}
		b.redis.Delete(ctx, "unavailable:ticket:"+ticketId)

	}
	return nil
}
