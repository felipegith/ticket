package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"ticket/ticket/domain/repositories"
)

type bookingRepository struct {
	dbContext *sql.DB
}

var _ repositories.BookingRepository = (*bookingRepository)(nil)

func NewBookingRepository(dbContext *sql.DB) repositories.BookingRepository {
	return &bookingRepository{
		dbContext: dbContext,
	}
}

func (repository *bookingRepository) Create(ctx context.Context, ticketId, seat, userId string) error {
	const query = `
		INSERT INTO bookings (id, ticket_id, seat, user_id)
		VALUES ($1, $2, $3, $4)
	`

	_, err := repository.dbContext.ExecContext(ctx, query, uuid.NewString(), ticketId, seat, userId)
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil
		}
		return fmt.Errorf("erro ao inserir booking: %w", err)
	}
	return nil
}
