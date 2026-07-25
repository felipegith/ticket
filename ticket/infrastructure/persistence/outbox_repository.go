package persistence

import (
	"context"
	"database/sql"

	"ticket/ticket/domain/entities"
	"ticket/ticket/domain/repositories"
)

type outboxRepository struct {
	dbContext *sql.DB
}

var _ repositories.OutboxRepository = (*outboxRepository)(nil)

func NewOutboxRepository(dbContext *sql.DB) repositories.OutboxRepository {
	return &outboxRepository{
		dbContext: dbContext,
	}
}

func (repository *outboxRepository) FetchPending(ctx context.Context, limit int) ([]*entities.Outbox, error) {
	const query = `
		SELECT id, type, content, created_at, updated_at
		FROM outbox
		WHERE updated_at IS NULL
		ORDER BY created_at
		LIMIT $1
	`

	rows, err := repository.dbContext.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*entities.Outbox, 0)
	for rows.Next() {
		message := &entities.Outbox{}

		var updatedAt sql.NullTime
		if err := rows.Scan(
			&message.Id,
			&message.Type,
			&message.Content,
			&message.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			message.UpdatedAt = updatedAt.Time
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (repository *outboxRepository) MarkAsProcessed(ctx context.Context, message *entities.Outbox) error {
	message.MarkAsSuccess()

	const query = `UPDATE outbox SET updated_at = $1 WHERE id = $2`

	_, err := repository.dbContext.ExecContext(ctx, query, message.UpdatedAt, message.Id)
	return err
}
