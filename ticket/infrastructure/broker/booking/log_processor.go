package booking

import (
	"context"
	"log"

	usecase "ticket/ticket/application/use-cases/booking"
)

type LogProcessor struct{}

func NewLogProcessor() *LogProcessor {
	return &LogProcessor{}
}

func (p *LogProcessor) Process(ctx context.Context, message *usecase.BookingInputModel) error {
	log.Printf("[booking] processando pagamento: ticketId=%s seat=%s userId=%s",
		message.TicketId, message.Seat, message.UserId)
	return nil
}
