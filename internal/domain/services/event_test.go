package services

import (
	"context"
	"testing"
	"time"

	"github.com/orensimple/otus_events_api/internal/domain/errors"
	"github.com/orensimple/otus_events_api/internal/domain/services"
	"github.com/orensimple/otus_events_api/internal/memory"
	"github.com/stretchr/testify/assert"
)

func TestGetEvents(t *testing.T) {
	eventStorage := memory.NewMemEventStorage()

	eventService := &services.EventService{
		EventStorage: eventStorage,
	}
	var i int64
	for i = 0; i < 10; i++ {
		time := time.Now()
		_, err := eventService.CreateEvent(context.Background(), i, "a", "a", "a", time, time)
		assert.NoError(t, err)
	}

	returnedEvents, err := eventService.GetEvents(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, 10, len(returnedEvents))
}

func TestDeleteEvent(t *testing.T) {
	eventStorage := memory.NewMemEventStorage()

	eventService := &services.EventService{
		EventStorage: eventStorage,
	}
	var i int64
	for i = 0; i < 2; i++ {
		time := time.Now()
		_, err := eventService.CreateEvent(context.Background(), i, "a", "a", "a", time, time)
		assert.NoError(t, err)
	}

	returnedEvents, err := eventService.GetEvents(context.Background())
	assert.NoError(t, err)

	err = eventService.DeleteEvent(context.Background(), 0)
	assert.NoError(t, err)

	err = eventService.DeleteEvent(context.Background(), 1)
	assert.NoError(t, err)

	err = eventService.DeleteEvent(context.Background(), 2)
	assert.Error(t, err)
	assert.Equal(t, err, errors.ErrEventNotFound)

	assert.Equal(t, 2, len(returnedEvents))
}

func TestCreateEvents(t *testing.T) {
	eventStorage := memory.NewMemEventStorage()

	eventService := &services.EventService{
		EventStorage: eventStorage,
	}

	time := time.Now()
	_, err := eventService.CreateEvent(context.Background(), 1, "a", "a", "a", time, time)
	assert.NoError(t, err)

	returnedEvents, err := eventService.GetEvents(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(returnedEvents))

	_, err = eventService.CreateEvent(context.Background(), 1, "a", "a", "a", time, time)
	assert.Error(t, err)
	assert.Equal(t, err, errors.ErrEventExist)
}

func TestUpdateEvents(t *testing.T) {
	eventStorage := memory.NewMemEventStorage()

	eventService := &services.EventService{
		EventStorage: eventStorage,
	}

	time := time.Now()
	_, err := eventService.CreateEvent(context.Background(), 1, "a", "a", "a", time, time)
	assert.NoError(t, err)

	returnedEvents, err := eventService.GetEvents(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(returnedEvents))

	_, err = eventService.UpdateEvent(context.Background(), 3, "a", "a", "a", time, time)
	assert.Error(t, err)
	assert.Equal(t, err, errors.ErrEventNotFound)

	_, err = eventService.UpdateEvent(context.Background(), 1, "b", "b", "a", time, time)
	assert.NoError(t, err)
}
