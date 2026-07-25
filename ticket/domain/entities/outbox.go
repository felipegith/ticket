package entities

import (
	"time"

	"github.com/google/uuid"
)

type Outbox struct {
	Id        string
	Type      string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOutbox(content string, eventType string) *Outbox {
	id := uuid.NewString()
	return &Outbox{
		Id:        id,
		Type:      eventType,
		Content:   content,
		CreatedAt: time.Now(),
	}
}

func (o *Outbox) MarkAsSuccess() {
	o.UpdatedAt = time.Now()
}
