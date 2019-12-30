package services

import (
	"context"
	"time"

	"github.com/orensimple/otus_events_api/internal/domain/interfaces"
	"github.com/orensimple/otus_events_api/internal/domain/models"
)

//EventService struct
type EventService struct {
	EventStorage interfaces.EventStorage
}

//CreateEvent func
func (es *EventService) CreateEvent(ctx context.Context, ID int64, owner, title, text string, startTime time.Time, endTime time.Time) (*models.Event, error) {
	event := &models.Event{
		ID:        ID,
		Owner:     owner,
		Title:     title,
		Text:      text,
		StartTime: startTime,
		EndTime:   endTime,
	}

	err := es.EventStorage.SaveEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

//UpdateEvent func
func (es *EventService) UpdateEvent(ctx context.Context, ID int64, owner, title, text string, startTime time.Time, endTime time.Time) (*models.Event, error) {
	event := &models.Event{
		ID:        ID,
		Owner:     owner,
		Title:     title,
		Text:      text,
		StartTime: startTime,
		EndTime:   endTime,
	}

	event, err := es.EventStorage.UpdateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

//GetEvents func
func (es *EventService) GetEvents(ctx context.Context) ([]*models.Event, error) {
	events, err := es.EventStorage.GetEvents(ctx)
	if err != nil {
		return nil, err
	}
	return events, nil
}

//GetEventsByTime func
func (es *EventService) GetEventsByTime(ctx context.Context, timeType string) ([]*models.Event, error) {
	switch timeType {
	case "day":
		events, err := es.EventStorage.GetEventsByDay(ctx)
		if err != nil {
			return nil, err
		}
		return events, nil
	case "week":
		events, err := es.EventStorage.GetEventsByWeek(ctx)
		if err != nil {
			return nil, err
		}
		return events, nil
	case "month":
		events, err := es.EventStorage.GetEventsByMonth(ctx)
		if err != nil {
			return nil, err
		}
		return events, nil
	default:
		events, err := es.EventStorage.GetEventsByDay(ctx)
		if err != nil {
			return nil, err
		}
		return events, nil
	}
}

//DeleteEvent func
func (es *EventService) DeleteEvent(ctx context.Context, ID int64) error {
	err := es.EventStorage.DeleteEvent(ctx, ID)

	return err
}
