package worker

import (
	"context"
	"log"
	"time"

	"ticket/ticket/domain/entities"
	"ticket/ticket/domain/repositories"
)

type Handler interface {
	Handle(ctx context.Context, message *entities.Outbox) error
}

type OutboxWorker struct {
	repository repositories.OutboxRepository
	handler    Handler
	interval   time.Duration
	batchSize  int
}

func NewOutboxWorker(
	repository repositories.OutboxRepository,
	handler Handler,
	interval time.Duration,
	batchSize int,
) *OutboxWorker {
	return &OutboxWorker{
		repository: repository,
		handler:    handler,
		interval:   interval,
		batchSize:  batchSize,
	}
}

func (w *OutboxWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("outbox worker iniciado (intervalo=%s, lote=%d)", w.interval, w.batchSize)

	for {
		select {
		case <-ctx.Done():
			log.Println("outbox worker encerrado")
			return
		case <-ticker.C:
			if err := w.processBatch(ctx); err != nil {
				log.Printf("outbox worker: erro ao processar lote: %v", err)
			}
		}
	}
}

func (w *OutboxWorker) processBatch(ctx context.Context) error {
	messages, err := w.repository.FetchPending(ctx, w.batchSize)
	if err != nil {
		return err
	}

	for _, message := range messages {

		if err := w.handler.Handle(ctx, message); err != nil {
			log.Printf("outbox worker: falha ao entregar mensagem %s: %v", message.Id, err)
			continue
		}

		if err := w.repository.MarkAsProcessed(ctx, message); err != nil {
			log.Printf("outbox worker: falha ao marcar mensagem %s como processada: %v", message.Id, err)
		}
	}

	return nil
}
