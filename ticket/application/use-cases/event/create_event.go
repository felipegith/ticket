package event

import (
	"context"
	"errors"
	"time"

	"ticket/ticket/domain/entities"
	"ticket/ticket/domain/repositories"
)

type EventInputModel struct {
	Name        string
	Description string
	Status      string
	Date        time.Time
}

type CreateEvent struct {
	repository repositories.EventRepository
}

func NewEvent(repository repositories.EventRepository) *CreateEvent {
	return &CreateEvent{
		repository: repository,
	}
}
func (r *CreateEvent) Execute(ctx context.Context, input EventInputModel) (string, error) {
	if err := ValidateEventInput(input); err != nil {
		return "", err
	}

	event := entities.NewEvent(input.Name, input.Description, input.Status, input.Date)

	if err := r.repository.Create(ctx, event); err != nil {
		return "", err
	}

	return event.Id, nil
}

func ValidateEventInput(input EventInputModel) error {
	if input.Name == "" || input.Description == "" || input.Status == "" {
		return errors.New("name, description, and status are required")
	}
	return nil
}
