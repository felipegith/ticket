package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"ticket/ticket/domain/entities"
	"ticket/ticket/domain/repositories"
)

type ticketRepository struct {
	dbContext *sql.DB
}

var _ repositories.TicketRepository = (*ticketRepository)(nil)

func NewTicketRepository(dbContext *sql.DB) repositories.TicketRepository {
	return &ticketRepository{
		dbContext: dbContext,
	}
}

func (repository *ticketRepository) Create(ctx context.Context, ticket *entities.Ticket) error {
	const insertTicket = `
		INSERT INTO tickets (id, event_id, price, user_id, status, seat)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := repository.dbContext.ExecContext(ctx, insertTicket, ticket.Id, ticket.EventId, ticket.Price, ticket.UserId, ticket.Status, ticket.Seat)
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entities.ErrSeatTaken
		}
		return fmt.Errorf("failed to insert ticket: %w", err)
	}
	return nil
}

func (repository *ticketRepository) SeatExists(ctx context.Context, eventID, seat string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM tickets WHERE event_id = $1 AND seat = $2)`

	var exists bool
	if err := repository.dbContext.QueryRowContext(ctx, query, eventID, seat).Scan(&exists); err != nil {
		return false, fmt.Errorf("erro ao verificar assento: %w", err)
	}
	return exists, nil
}

func (repository *ticketRepository) Update(ctx context.Context, ticketId string) (bool, error) {
	const query = `UPDATE tickets SET status = $1 WHERE id = $2`

	result, err := repository.dbContext.ExecContext(ctx, query, "sold", ticketId)
	if err != nil {
		return false, fmt.Errorf("erro ao atualizar ticket: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("erro ao obter linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
