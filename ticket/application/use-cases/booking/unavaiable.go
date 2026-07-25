package booking

import (
	"context"
	"encoding/json"
	"errors"
	"ticket/ticket/application/ports"
	"time"
)

type UnavailableInputModel struct {
	TicketId string
	Seat     string
	UserId   string
}

type Unavailable struct {
	redis ports.Cache
}

func NewUnavailable(redis ports.Cache) *Unavailable {
	return &Unavailable{
		redis: redis,
	}
}

func (c *Unavailable) Execute(ctx context.Context, input UnavailableInputModel) error {
	if err := ValidateInput(input); err != nil {
		return err
	}

	if payload, err := json.Marshal(input); err == nil {
		key := "unavailable:ticket:" + input.TicketId
		_ = c.redis.Set(ctx, key, string(payload), 6*time.Minute)
	}
	return nil
}

func ValidateInput(input UnavailableInputModel) error {
	if input.TicketId == "" || input.Seat == "" || input.UserId == "" {
		return errors.New("TicketId, userId and seat are required")
	}
	return nil
}
