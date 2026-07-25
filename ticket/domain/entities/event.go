package entities

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id          string
	Name        string
	Description string
	Status      string
	Date        time.Time
	Tickets     []Ticket
}

func NewEvent(name string, description string, status string, date time.Time) *Event {
	id := uuid.NewString()
	return &Event{
		Id:          id,
		Name:        name,
		Description: description,
		Status:      status,
		Date:        date,
	}
}
