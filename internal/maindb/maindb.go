package maindb

import (
	"context"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/orensimple/otus_events_api/internal/domain/models"
	"github.com/orensimple/otus_events_api/internal/logger"
)

// implements domain.interfaces.EventStorage
type PgEventStorage struct {
	db *sqlx.DB
}

func NewPgEventStorage(dsn string) (*PgEventStorage, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PgEventStorage{db: db}, nil
}

func (pges *PgEventStorage) SaveEvent(ctx context.Context, event *models.Event) error {
	query := `
		INSERT INTO events(id, owner, title, text, start_time, end_time)
		VALUES (:id, :owner, :title, :text, :start_time, :end_time)
	`
	_, err := pges.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":         event.ID,
		"owner":      event.Owner,
		"title":      event.Title,
		"text":       event.Text,
		"start_time": event.StartTime,
		"end_time":   event.EndTime,
	})
	return err
}

// UpdateEvent to maindb
func (pges *PgEventStorage) UpdateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	query := `
		UPDATE events
		SET owner = :owner, title = :title, text = :text, start_time = :start_time, end_time =:end_time
		WHERE id = :id
	`
	_, err := pges.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":         event.ID,
		"owner":      event.Owner,
		"title":      event.Title,
		"text":       event.Text,
		"start_time": event.StartTime,
		"end_time":   event.EndTime,
	})
	return event, err
}

// GetEvents from maindb
func (pges *PgEventStorage) GetEvents(ctx context.Context) ([]*models.Event, error) {

	return nil, nil
}

// GetEventsByDay from maindb
func (pges *PgEventStorage) GetEventsByDay(ctx context.Context) ([]*models.Event, error) {
	query := `
	SELECT id, owner, title, text, start_time, end_time FROM events
	WHERE start_time BETWEEN
	date_trunc('day', now())
	AND date_trunc('day', now()) +('1 day')::interval -('1 second')::interval
`
	rows, err := pges.db.QueryContext(ctx, query)
	if err != nil {
		logger.ContextLogger.Infof("QueryContext", "err", err.Error())
	}
	defer rows.Close()
	var events []*models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Owner, &event.Title, &event.Text, &event.StartTime, &event.EndTime); err != nil {
			logger.ContextLogger.Infof("rowScan", "err", err.Error())
		}
		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		logger.ContextLogger.Infof("rowErr", "err", err.Error())
	}

	return events, nil
}

// GetEventsByWeek from maindb
func (pges *PgEventStorage) GetEventsByWeek(ctx context.Context) ([]*models.Event, error) {
	query := `
	SELECT id, owner, title, text, start_time, end_time FROM events
	WHERE start_time BETWEEN
	date_trunc('week', now())
	AND date_trunc('week', now()) +('1 week')::interval -('1 day')::interval
`
	rows, err := pges.db.QueryContext(ctx, query)
	if err != nil {
		logger.ContextLogger.Infof("QueryContext", "err", err.Error())
	}
	defer rows.Close()
	var events []*models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Owner, &event.Title, &event.Text, &event.StartTime, &event.EndTime); err != nil {
			logger.ContextLogger.Infof("rowScan", "err", err.Error())
		}
		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		logger.ContextLogger.Infof("rowErr", "err", err.Error())
	}

	return events, nil
}

// GetEventsByMonth from maindb
func (pges *PgEventStorage) GetEventsByMonth(ctx context.Context) ([]*models.Event, error) {
	query := `
	SELECT id, owner, title, text, start_time, end_time FROM events
	WHERE start_time BETWEEN
	date_trunc('month', now())
	AND date_trunc('month', now()) +('1 month')::interval -('1 day')::interval
`
	rows, err := pges.db.QueryContext(ctx, query)
	if err != nil {
		logger.ContextLogger.Infof("QueryContext", "err", err.Error())
	}
	defer rows.Close()
	var events []*models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Owner, &event.Title, &event.Text, &event.StartTime, &event.EndTime); err != nil {
			logger.ContextLogger.Infof("rowScan", "err", err.Error())
		}
		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		logger.ContextLogger.Infof("rowErr", "err", err.Error())
	}

	return events, nil
}

// DeleteEvent from maindb
func (pges *PgEventStorage) DeleteEvent(ctx context.Context, ID int64) error {
	query := `
		DELETE FROM events
		WHERE id = :id
	`
	_, err := pges.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id": ID,
	})
	return err
}
