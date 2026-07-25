package event

import (
	"context"
	"encoding/json"
	"time"

	"ticket/ticket/application/ports"
	"ticket/ticket/domain/repositories"
)

const eventsCacheKey = "events:all"

type TicketOutputModel struct {
	Id     string  `json:"id"`
	Status string  `json:"status"`
	Seat   string  `json:"seat"`
	Price  float64 `json:"price"`
}

type EventOutputModel struct {
	Id          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Status      string              `json:"status"`
	Date        time.Time           `json:"date"`
	Tickets     []TicketOutputModel `json:"tickets"`
}

type GetEvents struct {
	repository repositories.EventRepository
	cache      ports.Cache
	ttl        time.Duration
}

func NewGetEvents(repository repositories.EventRepository, cache ports.Cache, ttl time.Duration) *GetEvents {
	return &GetEvents{
		repository: repository,
		cache:      cache,
		ttl:        ttl,
	}
}

func (r *GetEvents) Execute(ctx context.Context) ([]EventOutputModel, error) {
	if cached, found, err := r.cache.Get(ctx, eventsCacheKey); err == nil && found {
		var events []EventOutputModel
		if err := json.Unmarshal([]byte(cached), &events); err == nil {
			return events, nil
		}
	}

	events, err := r.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	outputEvents := make([]EventOutputModel, len(events))
	for i, event := range events {
		tickets := make([]TicketOutputModel, len(event.Tickets))
		for j, t := range event.Tickets {
			tickets[j] = TicketOutputModel{
				Status: t.Status,
				Seat:   t.Seat,
				Price:  t.Price,
			}
		}
		outputEvents[i] = EventOutputModel{
			Id:          event.Id,
			Name:        event.Name,
			Description: event.Description,
			Status:      event.Status,
			Date:        event.Date,
			Tickets:     tickets,
		}
	}

	if payload, err := json.Marshal(outputEvents); err == nil {
		_ = r.cache.Set(ctx, eventsCacheKey, string(payload), r.ttl)
	}

	return outputEvents, nil
}
