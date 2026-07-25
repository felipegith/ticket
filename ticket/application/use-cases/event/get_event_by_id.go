package event

import (
	"context"

	"ticket/ticket/application/ports"
	"ticket/ticket/domain/repositories"
)

const statusUnavailable = "unavailable"

type GetEventById struct {
	repository repositories.EventRepository
	redis      ports.Cache
}

func NewGetEventById(repository repositories.EventRepository, redis ports.Cache) *GetEventById {
	return &GetEventById{
		repository: repository,
		redis:      redis,
	}
}

func (g *GetEventById) Execute(ctx context.Context, id string) (*EventOutputModel, error) {
	event, err := g.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	tickets := make([]TicketOutputModel, len(event.Tickets))
	for i, t := range event.Tickets {
		status := t.Status

		if _, found, err := g.redis.Get(ctx, "unavailable:ticket:"+t.Id); err == nil && found {
			status = statusUnavailable
		}

		tickets[i] = TicketOutputModel{
			Id:     t.Id,
			Status: status,
			Seat:   t.Seat,
			Price:  t.Price,
		}
	}
	return &EventOutputModel{
		Id:          event.Id,
		Name:        event.Name,
		Description: event.Description,
		Status:      event.Status,
		Date:        event.Date,
		Tickets:     tickets,
	}, nil
}
