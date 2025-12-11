package models

import (
	"context"
	"database/sql"
	"sync"
	"time"
)

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
}

type EventRepository struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(ctx context.Context, event *Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	query := `INSERT INTO events (title, description, start_time, end_time, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, event.Title, event.Description, event.StartTime, event.EndTime).Scan(&event.ID, &event.CreatedAt)
}

func (r *EventRepository) GetByID(ctx context.Context, id string) (*Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event := &Event{}
	query := `SELECT id, title, description, start_time, end_time, created_at FROM events WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.CreatedAt)
	return event, err
}

func (r *EventRepository) GetAll(ctx context.Context) ([]Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query := `SELECT id, title, description, start_time, end_time, created_at FROM events`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (r *EventRepository) Update(ctx context.Context, event *Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	query := `UPDATE events SET title = $1, description = $2, start_time = $3, end_time = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, event.Title, event.Description, event.StartTime, event.EndTime, event.ID)
	return err
}

func (r *EventRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	query := `DELETE FROM events WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}