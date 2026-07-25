package ticket

import (
	"context"
	"errors"
	"ticket/ticket/domain/entities"
	"ticket/ticket/domain/repositories"
)

type TicketInputModel struct {
	EventId string
	Price   float64
	Status  string
	Seat    string
	UserId  string
}

type CreateTicket struct {
	repository repositories.TicketRepository
}

func NewTicket(repository repositories.TicketRepository) *CreateTicket {
	return &CreateTicket{
		repository: repository,
	}
}

func (c *CreateTicket) Execute(ctx context.Context, input TicketInputModel) (string, error) {
	if err := ValidateTicketInput(input); err != nil {
		return "", err
	}

	exists, err := c.repository.SeatExists(ctx, input.EventId, input.Seat)
	if err != nil {
		return "", err
	}

	if exists {
		return "", entities.ErrSeatTaken
	}

	ticket := entities.NewTicket(input.EventId, input.Price, input.Status, input.Seat, input.UserId)
	if err := c.repository.Create(ctx, ticket); err != nil {
		return "", err
	}
	return ticket.Id, nil
}

func ValidateTicketInput(input TicketInputModel) error {
	if input.EventId == "" || input.Price <= 0 || input.Status == "" || input.Seat == "" || input.UserId == "" {
		return errors.New("eventId, price, status, seat and userId are required")
	}
	return nil
}
