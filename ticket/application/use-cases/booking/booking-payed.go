package booking

import (
	"context"
)

type Handler interface {
	Handle(ctx context.Context, message *BookingInputModel) error
}
type BookingInputModel struct {
	TicketId string
	Seat     string
	UserId   string
}

type BookingPayed struct {
	handler Handler
}

func NewBookingPayed(handler Handler) *BookingPayed {
	return &BookingPayed{
		handler: handler,
	}
}

func (b *BookingPayed) Execute(ctx context.Context, input BookingInputModel) error {
	return b.handler.Handle(ctx, &input)
}
