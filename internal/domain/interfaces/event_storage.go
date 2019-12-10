package interfaces

import (
	"context"

	"github.com/orensimple/otus_events_api/internal/domain/models"
)

// EventStorage interface
type EventStorage interface {
	SaveEvent(ctx context.Context, event *models.Event) error
	UpdateEvent(ctx context.Context, event *models.Event) (*models.Event, error)
	GetEvents(ctx context.Context) ([]*models.Event, error)
	GetEventsByDay(ctx context.Context) ([]*models.Event, error)
	GetEventsByWeek(ctx context.Context) ([]*models.Event, error)
	GetEventsByMonth(ctx context.Context) ([]*models.Event, error)
	DeleteEvent(ctx context.Context, ID int64) error
}
