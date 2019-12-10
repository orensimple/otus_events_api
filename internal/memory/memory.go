package memory

import (
	"context"
	"sync"
	"time"

	"github.com/orensimple/otus_events_api/internal/domain/errors"
	"github.com/orensimple/otus_events_api/internal/domain/models"
)

// MemEventStorage slice
type MemEventStorage struct {
	events map[int64]*models.Event
	mutex  *sync.Mutex
}

// NewMemEventStorage init
func NewMemEventStorage() *MemEventStorage {
	return &MemEventStorage{
		events: make(map[int64]*models.Event),
		mutex:  new(sync.Mutex),
	}
}

// SaveEvent to mem
func (mem *MemEventStorage) SaveEvent(ctx context.Context, event *models.Event) error {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()
	if _, ok := mem.events[event.ID]; ok {
		return errors.ErrEventExist

	}
	mem.events[event.ID] = event
	return nil
}

// UpdateEvent to mem
func (mem *MemEventStorage) UpdateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	if _, ok := mem.events[event.ID]; ok {
		mem.mutex.Lock()
		mem.events[event.ID] = event
		mem.mutex.Unlock()

		return event, nil
	}

	return nil, errors.ErrEventNotFound
}

// GetEvents from mem
func (mem *MemEventStorage) GetEvents(ctx context.Context) ([]*models.Event, error) {
	Events := make([]*models.Event, 0)
	mem.mutex.Lock()
	for _, bm := range mem.events {
		Events = append(Events, bm)
	}
	mem.mutex.Unlock()

	return Events, nil
}

// GetEventsByDay from mem
func (mem *MemEventStorage) GetEventsByDay(ctx context.Context) ([]*models.Event, error) {
	timeNow := time.Now()
	Events := make([]*models.Event, 0)
	mem.mutex.Lock()
	for _, bm := range mem.events {
		if timeNow.YearDay() == bm.StartTime.YearDay() {
			Events = append(Events, bm)
		}
	}
	mem.mutex.Unlock()

	return Events, nil
}

// GetEventsByWeek from mem
func (mem *MemEventStorage) GetEventsByWeek(ctx context.Context) ([]*models.Event, error) {
	timeNow := time.Now()
	weekNow, yearNow := timeNow.ISOWeek()
	Events := make([]*models.Event, 0)
	mem.mutex.Lock()
	for _, bm := range mem.events {
		weekEvent, yearEvent := bm.StartTime.ISOWeek()
		if weekNow == weekEvent && yearNow == yearEvent {
			Events = append(Events, bm)
		}
	}
	mem.mutex.Unlock()

	return Events, nil
}

// GetEventsByMonth from mem
func (mem *MemEventStorage) GetEventsByMonth(ctx context.Context) ([]*models.Event, error) {
	timeNow := time.Now()
	Events := make([]*models.Event, 0)
	mem.mutex.Lock()
	for _, bm := range mem.events {
		if timeNow.Month() == bm.StartTime.Month() && timeNow.Year() == bm.StartTime.Year() {
			Events = append(Events, bm)
		}
	}
	mem.mutex.Unlock()

	return Events, nil
}

// DeleteEvent from mem
func (mem *MemEventStorage) DeleteEvent(ctx context.Context, ID int64) error {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()

	_, ex := mem.events[ID]
	if ex {
		delete(mem.events, ID)
		return nil
	}

	return errors.ErrEventNotFound
}
