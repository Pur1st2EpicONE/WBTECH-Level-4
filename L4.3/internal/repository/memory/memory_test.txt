package memory

import (
	"testing"
	"time"

	"L4.3/internal/config"
	"L4.3/internal/models"
	"L4.3/pkg/logger/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStorage_CreateEvent(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mocks.NewMockLogger(controller)
	mockLogger.EXPECT().Debug("repository — new user created", "UserID", 42, "layer", "repository.memory").Times(1)

	storage := NewStorage(config.Storage{ExpectedUsers: 1, MaxEventsPerUser: 1, MaxEventsPerDay: 1}, mockLogger)
	eventDate := time.Date(2025, 12, 3, 10, 0, 0, 0, time.UTC)

	event := &models.Event{
		Meta: models.Meta{UserID: 42, EventDate: eventDate},
		Data: models.Data{Text: "aboba"},
	}

	id, err := storage.CreateEvent(event)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	found := storage.GetEventByID(id)
	require.NotNil(t, found)
	require.Equal(t, "aboba", found.Data.Text)

}

func TestStorage_UpdateEvent(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mocks.NewMockLogger(controller)
	mockLogger.EXPECT().Debug("repository — new user created", "UserID", 7, "layer", "repository.memory").Times(1)
	mockLogger.EXPECT().Debug("repository — event data updated", "UserID", 7, "EventID", gomock.Any(), "layer", "repository.memory").Times(1)
	mockLogger.EXPECT().Debug("repository — event meta updated", "UserID", 7, "EventID", gomock.Any(), "layer", "repository.memory").Times(1)

	storage := NewStorage(config.Storage{}, mockLogger)

	eventDate := time.Date(2025, 12, 1, 9, 0, 0, 0, time.UTC)
	event := &models.Event{
		Meta: models.Meta{UserID: 7, EventDate: eventDate},
		Data: models.Data{Text: "old, not cool text"},
	}

	id, err := storage.CreateEvent(event)
	require.NoError(t, err)

	updatedEvent := &models.Event{
		Meta: models.Meta{EventID: id, UserID: 7, NewDate: eventDate.Add(24 * time.Hour)},
		Data: models.Data{Text: "new, really cool text"},
	}

	err = storage.UpdateEvent(updatedEvent)
	require.NoError(t, err)

	oldMeta := &models.Meta{UserID: 7, EventDate: eventDate}
	oldEvents, err := storage.GetEvents(oldMeta, models.Day)
	require.NoError(t, err)
	require.Len(t, oldEvents, 0)

	newMeta := &models.Meta{UserID: 7, EventDate: eventDate.Add(24 * time.Hour)}
	newEvents, err := storage.GetEvents(newMeta, models.Day)
	require.NoError(t, err)
	require.Len(t, newEvents, 1)
	require.Equal(t, "new, really cool text", newEvents[0].Data.Text)

	found := storage.GetEventByID(id)
	require.NotNil(t, found)
	require.Equal(t, "new, really cool text", found.Data.Text)
	require.True(t, found.Meta.EventDate.Equal(eventDate.Add(24*time.Hour)))

	count, err := storage.CountUserEvents(7)
	require.NoError(t, err)
	require.Equal(t, 1, count)

}

func TestStorage_UpdateEvent_Else(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mocks.NewMockLogger(controller)
	mockLogger.EXPECT().Debug("repository — new user created", "UserID", 7, "layer", "repository.memory").Times(1)
	mockLogger.EXPECT().Debug("repository — event meta updated", "UserID", 7, "EventID", gomock.Any(), "layer", "repository.memory").Times(1)

	storage := NewStorage(config.Storage{}, mockLogger)

	eventDate := time.Date(2025, 12, 1, 9, 0, 0, 0, time.UTC)

	event1 := &models.Event{
		Meta: models.Meta{UserID: 7, EventDate: eventDate},
		Data: models.Data{Text: "first"},
	}

	event2 := &models.Event{
		Meta: models.Meta{UserID: 7, EventDate: eventDate},
		Data: models.Data{Text: "second"},
	}

	id1, err := storage.CreateEvent(event1)
	require.NoError(t, err)
	_, err = storage.CreateEvent(event2)
	require.NoError(t, err)

	updatedEvent := &models.Event{
		Meta: models.Meta{EventID: id1, UserID: 7, NewDate: eventDate.Add(24 * time.Hour)},
		Data: models.Data{Text: "first"},
	}

	err = storage.UpdateEvent(updatedEvent)
	require.NoError(t, err)

	oldMeta := &models.Meta{UserID: 7, EventDate: eventDate}
	remainingEvents, err := storage.GetEvents(oldMeta, models.Day)
	require.NoError(t, err)
	require.Len(t, remainingEvents, 1)
	require.Equal(t, "second", remainingEvents[0].Data.Text)

}

func TestStorage_DeleteEvent_AllBranches(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mocks.NewMockLogger(controller)
	mockLogger.EXPECT().Debug("repository — new user created", "UserID", 5, "layer", "repository.memory").Times(1)
	mockLogger.EXPECT().Debug("repository — new user created", "UserID", 6, "layer", "repository.memory").Times(1)

	storage := NewStorage(config.Storage{}, mockLogger)

	eventDate1 := time.Date(2025, 12, 2, 8, 0, 0, 0, time.UTC)
	event1 := &models.Event{
		Meta: models.Meta{UserID: 5, EventDate: eventDate1},
		Data: models.Data{Text: "Qwe? Qwe!"},
	}

	id1, err := storage.CreateEvent(event1)
	require.NoError(t, err)

	err = storage.DeleteEvent(&models.Meta{EventID: id1})
	require.NoError(t, err)

	found := storage.GetEventByID(id1)
	require.Nil(t, found)

	count, err := storage.CountUserEvents(5)
	require.NoError(t, err)
	require.Equal(t, 0, count)

	meta1 := &models.Meta{UserID: 5, EventDate: eventDate1}
	events1, err := storage.GetEvents(meta1, models.Day)
	require.NoError(t, err)
	require.Len(t, events1, 0)

	eventDate2 := time.Date(2025, 12, 2, 8, 0, 0, 0, time.UTC)

	event2a := &models.Event{Meta: models.Meta{UserID: 6, EventDate: eventDate2}, Data: models.Data{Text: "first"}}
	event2b := &models.Event{Meta: models.Meta{UserID: 6, EventDate: eventDate2}, Data: models.Data{Text: "second"}}

	id2a, err := storage.CreateEvent(event2a)
	require.NoError(t, err)
	_, err = storage.CreateEvent(event2b)
	require.NoError(t, err)

	err = storage.DeleteEvent(&models.Meta{EventID: id2a})
	require.NoError(t, err)

	events2, err := storage.GetEvents(&models.Meta{UserID: 6, EventDate: eventDate2}, models.Day)
	require.NoError(t, err)
	require.Len(t, events2, 1)
	require.Equal(t, "second", events2[0].Data.Text)

}

func TestStorage_GetEvents_Periods(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mocks.NewMockLogger(controller)
	mockLogger.EXPECT().Debug("repository — new user created", "UserID", 8, "layer", "repository.memory").Times(1)

	storage := NewStorage(config.Storage{}, mockLogger)

	baseDate := time.Date(2025, 12, 1, 12, 0, 0, 0, time.UTC)

	_, err := storage.CreateEvent(&models.Event{
		Meta: models.Meta{UserID: 8, EventDate: baseDate},
		Data: models.Data{Text: "average event"},
	})
	require.NoError(t, err)

	_, err = storage.CreateEvent(&models.Event{
		Meta: models.Meta{UserID: 8, EventDate: baseDate.Add(24 * time.Hour)},
		Data: models.Data{Text: "lame event"},
	})
	require.NoError(t, err)

	_, err = storage.CreateEvent(&models.Event{
		Meta: models.Meta{UserID: 8, EventDate: baseDate.Add(-24 * time.Hour)},
		Data: models.Data{Text: "cool event"},
	})
	require.NoError(t, err)

	metaDay := &models.Meta{UserID: 8, EventDate: baseDate}
	dayEvents, err := storage.GetEvents(metaDay, models.Day)
	require.NoError(t, err)
	require.Len(t, dayEvents, 1)

	metaWeek := &models.Meta{UserID: 8, EventDate: baseDate}
	weekEvents, err := storage.GetEvents(metaWeek, models.Week)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(weekEvents), 2)

	metaMonth := &models.Meta{UserID: 8, EventDate: baseDate}
	monthEvents, err := storage.GetEvents(metaMonth, models.Month)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(monthEvents), 2)

	_, err = storage.GetEvents(metaMonth, "hour")
	require.Error(t, err)

	metaNotFound := &models.Meta{UserID: 99, EventDate: time.Now()}
	events, err := storage.GetEvents(metaNotFound, models.Day)
	require.NoError(t, err)
	require.Empty(t, events)

}

func TestStorage_Close(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mocks.NewMockLogger(controller)

	mockLogger.EXPECT().Debug("repository — new user created", "UserID", 11, "layer", "repository.memory").Times(1)
	mockLogger.EXPECT().LogInfo("in-memory storage — cleared and stopped", "layer", "repository.memory").Times(1)

	storage := NewStorage(config.Storage{}, mockLogger)

	event := &models.Event{
		Meta: models.Meta{UserID: 11, EventDate: time.Now()},
		Data: models.Data{Text: "hello there"},
	}

	_, err := storage.CreateEvent(event)
	require.NoError(t, err)

	storage.Close()

	require.Nil(t, storage.db)
	require.Nil(t, storage.eventsByID)
	require.Nil(t, storage.userEventCount)

}
