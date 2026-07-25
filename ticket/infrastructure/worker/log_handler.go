package worker

import (
	"context"
	"log"

	"ticket/ticket/domain/entities"
)

type LogHandler struct{}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (h *LogHandler) Handle(ctx context.Context, message *entities.Outbox) error {
	log.Printf("[handler] processando outbox id=%s type=%s content=%s",
		message.Id, message.Type, message.Content)
	return nil
}
