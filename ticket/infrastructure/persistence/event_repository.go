package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"ticket/ticket/domain/entities"
	"ticket/ticket/domain/repositories"
)

type eventRepository struct {
	dbContext *sql.DB
}

var _ repositories.EventRepository = (*eventRepository)(nil)

func NewEventRepository(dbContext *sql.DB) repositories.EventRepository {
	return &eventRepository{
		dbContext: dbContext,
	}
}

func (repository *eventRepository) Create(ctx context.Context, event *entities.Event) error {
	transaction, err := repository.dbContext.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	const insertEvent = `
		INSERT INTO events (id, name, description, status, date)
		VALUES ($1, $2, $3, $4, $5)
	`
	if _, err := transaction.ExecContext(ctx, insertEvent,
		event.Id,
		event.Name,
		event.Description,
		event.Status,
		event.Date,
	); err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	content, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event for outbox: %w", err)
	}

	outbox := entities.NewOutbox(string(content), "EventCreated")

	const insertOutbox = `
		INSERT INTO outbox (id, type, content, created_at)
		VALUES ($1, $2, $3, $4)
	`
	if _, err := transaction.ExecContext(ctx, insertOutbox,
		outbox.Id,
		outbox.Type,
		outbox.Content,
		outbox.CreatedAt,
	); err != nil {
		return fmt.Errorf("failed to insert outbox entry: %w", err)
	}

	return transaction.Commit()
}

func (repository *eventRepository) GetAll(ctx context.Context) ([]*entities.Event, error) {

	const query = `
		SELECT
			events.id,
			events.name,
			events.description,
			events.status,
			events.date,
			tickets.status,
			tickets.seat,
			tickets.price
		FROM events
		LEFT JOIN tickets ON events.id = tickets.event_id
		ORDER BY events.id
	`

	rows, err := repository.dbContext.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventsByID := make(map[string]*entities.Event)
	events := make([]*entities.Event, 0)

	for rows.Next() {
		var (
			id, name, description, status string
			date                          time.Time
			ticketStatus, ticketSeat      sql.NullString
			ticketPrice                   sql.NullFloat64
		)
		if err := rows.Scan(
			&id, &name, &description, &status, &date,
			&ticketStatus, &ticketSeat, &ticketPrice,
		); err != nil {
			return nil, err
		}

		event, ok := eventsByID[id]
		if !ok {
			event = &entities.Event{
				Id:          id,
				Name:        name,
				Description: description,
				Status:      status,
				Date:        date,
			}
			eventsByID[id] = event
			events = append(events, event)
		}

		if ticketStatus.Valid {
			event.Tickets = append(event.Tickets, entities.Ticket{
				EventId: id,
				Status:  ticketStatus.String,
				Seat:    ticketSeat.String,
				Price:   ticketPrice.Float64,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (repository *eventRepository) GetByID(ctx context.Context, id string) (*entities.Event, error) {

	const query = `
		SELECT
			events.id,
			events.name,
			events.description,
			events.status,
			events.date,
			tickets.id,
			tickets.status,
			tickets.seat,
			tickets.price
		FROM events
		LEFT JOIN tickets ON events.id = tickets.event_id
		WHERE events.id = $1
	`

	rows, err := repository.dbContext.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var event *entities.Event
	for rows.Next() {
		var (
			eventID, name, description, status string
			date                               time.Time
			ticketID, ticketStatus, ticketSeat sql.NullString
			ticketPrice                        sql.NullFloat64
		)
		if err := rows.Scan(
			&eventID, &name, &description, &status, &date,
			&ticketID, &ticketStatus, &ticketSeat, &ticketPrice,
		); err != nil {
			return nil, err
		}

		if event == nil {
			event = &entities.Event{
				Id:          eventID,
				Name:        name,
				Description: description,
				Status:      status,
				Date:        date,
			}
		}

		if ticketID.Valid {
			event.Tickets = append(event.Tickets, entities.Ticket{
				Id:      ticketID.String,
				EventId: eventID,
				Status:  ticketStatus.String,
				Seat:    ticketSeat.String,
				Price:   ticketPrice.Float64,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if event == nil {
		return nil, entities.ErrEventNotFound
	}

	return event, nil
}
