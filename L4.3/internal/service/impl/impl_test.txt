package impl

import (
	"testing"
	"time"

	"L4.3/internal/config"
	"L4.3/internal/errs"

	"L4.3/internal/models"
	storageMock "L4.3/internal/repository/mocks"
	loggerMock "L4.3/pkg/logger/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateEvent_Success(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	now := time.Now().UTC().Add(24 * time.Hour)
	event := &models.Event{
		Meta: models.Meta{UserID: 1, EventDate: now},
		Data: models.Data{Text: "ok"},
	}

	mockLogger.EXPECT().Debug("service — user 1 has 5 remaining event slots", "UserID", event.Meta.UserID, "layer", "service.impl")
	mockStorage.EXPECT().CountUserEvents(event.Meta.UserID).Return(0, nil)
	mockStorage.EXPECT().CreateEvent(event).Return("result id", nil)

	id, err := service.CreateEvent(event)
	assert.NoError(t, err)
	assert.Equal(t, "result id", id)

}

func TestCreateEvent_ErrValidateCreate(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	event := &models.Event{
		Meta: models.Meta{UserID: 0, EventDate: time.Now().Add(24 * time.Hour)},
		Data: models.Data{Text: "ok"},
	}

	id, err := service.CreateEvent(event)
	assert.ErrorIs(t, err, errs.ErrInvalidUserID)
	assert.Equal(t, "", id)

}

func TestCreateEvent_ErrCountUserEvents(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)
	now := time.Now().Add(24 * time.Hour)

	event := &models.Event{
		Meta: models.Meta{UserID: 1, EventDate: now},
		Data: models.Data{Text: "ok"},
	}

	mockStorage.EXPECT().CountUserEvents(event.Meta.UserID).Return(0, assert.AnError)

	id, err := service.CreateEvent(event)
	assert.Equal(t, "", id)
	assert.ErrorIs(t, err, assert.AnError)

}

func TestCreateEvent_ErrMaxEvents(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)
	now := time.Now().Add(24 * time.Hour)

	event := &models.Event{
		Meta: models.Meta{UserID: 1, EventDate: now},
		Data: models.Data{Text: "ok"},
	}

	mockStorage.EXPECT().CountUserEvents(event.Meta.UserID).Return(5, nil)
	mockLogger.EXPECT().Debug("service — user 1 has 0 remaining event slots", "UserID", event.Meta.UserID, "layer", "service.impl")

	id, err := service.CreateEvent(event)
	assert.Equal(t, "", id)
	assert.ErrorIs(t, err, errs.ErrMaxEvents)

}

func TestUpdateEvent_Success(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	now := time.Now().UTC().Add(24 * time.Hour)
	eventID := uuid.New().String()

	oldEvent := &models.Event{
		Meta: models.Meta{UserID: 1, EventID: eventID, EventDate: now},
		Data: models.Data{Text: "old"},
	}

	event := &models.Event{
		Meta: models.Meta{UserID: 1, EventID: eventID, NewDate: now.Add(24 * time.Hour)},
		Data: models.Data{Text: "new"},
	}

	mockStorage.EXPECT().GetEventByID(eventID).Return(oldEvent)
	mockStorage.EXPECT().UpdateEvent(event).Return(nil)

	assert.NoError(t, service.UpdateEvent(event))

}

func TestUpdateEvent_ErrInvalidUserID(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	event := &models.Event{
		Meta: models.Meta{UserID: 0, EventID: uuid.New().String()},
	}

	assert.ErrorIs(t, service.UpdateEvent(event), errs.ErrInvalidUserID)

}

func TestUpdateEvent_ErrEventNotFound(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)
	eventID := uuid.New().String()

	event := &models.Event{
		Meta: models.Meta{UserID: 1, EventID: eventID},
	}

	mockStorage.EXPECT().GetEventByID(eventID).Return((*models.Event)(nil))

	assert.ErrorIs(t, service.UpdateEvent(event), errs.ErrEventNotFound)

}

func TestUpdateEvent_ErrUnauthorized(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	now := time.Now().UTC().Add(24 * time.Hour)
	eventID := uuid.New().String()

	oldEvent := &models.Event{
		Meta: models.Meta{UserID: 2, EventID: eventID, EventDate: now},
		Data: models.Data{Text: "old"},
	}

	event := &models.Event{
		Meta: models.Meta{UserID: 1, EventID: eventID, NewDate: now.Add(24 * time.Hour)},
		Data: models.Data{Text: "new"},
	}

	mockStorage.EXPECT().GetEventByID(eventID).Return(oldEvent)

	assert.ErrorIs(t, service.UpdateEvent(event), errs.ErrUnauthorized)

}

func TestUpdateEvent_ErrNothingToUpdate(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	now := time.Now().UTC().Add(24 * time.Hour)
	eventID := uuid.New().String()

	oldEvent := &models.Event{
		Meta: models.Meta{UserID: 1, EventID: eventID, EventDate: now},
		Data: models.Data{Text: "same"},
	}

	event := &models.Event{
		Meta: models.Meta{UserID: 1, EventID: eventID},
		Data: models.Data{Text: "same"},
	}

	mockStorage.EXPECT().GetEventByID(eventID).Return(oldEvent)

	assert.ErrorIs(t, service.UpdateEvent(event), errs.ErrNothingToUpdate)

}

func TestDeleteEvent_Success(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	now := time.Now().UTC().Add(24 * time.Hour)
	meta := &models.Meta{UserID: 1, EventID: uuid.New().String()}

	oldEvent := &models.Event{
		Meta: models.Meta{UserID: 1, EventID: meta.EventID, EventDate: now},
		Data: models.Data{Text: "ok"},
	}

	mockStorage.EXPECT().GetEventByID(meta.EventID).Return(oldEvent)
	mockStorage.EXPECT().DeleteEvent(meta).Return(nil)

	assert.NoError(t, service.DeleteEvent(meta))

}

func TestDeleteEvent_ErrInvalidUserID(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	meta := &models.Meta{UserID: 0, EventID: uuid.New().String()}

	assert.ErrorIs(t, service.DeleteEvent(meta), errs.ErrInvalidUserID)

}

func TestDeleteEvent_ErrEventNotFound(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	meta := &models.Meta{UserID: 1, EventID: uuid.New().String()}

	mockStorage.EXPECT().GetEventByID(meta.EventID).Return((*models.Event)(nil))

	assert.ErrorIs(t, service.DeleteEvent(meta), errs.ErrEventNotFound)

}

func TestDeleteEvent_ErrUnauthorized(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	meta := &models.Meta{UserID: 1, EventID: uuid.New().String()}
	oldEvent := &models.Event{Meta: models.Meta{UserID: 2, EventID: meta.EventID}}

	mockStorage.EXPECT().GetEventByID(meta.EventID).Return(oldEvent)

	assert.ErrorIs(t, service.DeleteEvent(meta), errs.ErrUnauthorized)

}

func TestGetEvents_SuccessSorted(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	meta := &models.Meta{UserID: 1, EventDate: time.Now().UTC().Add(24 * time.Hour)}
	soon := models.Event{Meta: models.Meta{UserID: 1, EventDate: meta.EventDate.Add(24 * time.Hour)}, Data: models.Data{Text: "soon"}}
	later := models.Event{Meta: models.Meta{UserID: 1, EventDate: meta.EventDate.Add(48 * time.Hour)}, Data: models.Data{Text: "later"}}
	earlier := models.Event{Meta: models.Meta{UserID: 1, EventDate: meta.EventDate.Add(0)}, Data: models.Data{Text: "earlier"}}

	unsorted := []models.Event{soon, later, earlier}

	mockStorage.EXPECT().GetEvents(meta, models.Day).Return(unsorted, nil)

	events, err := service.GetEvents(meta, models.Day)
	assert.NoError(t, err)

	if assert.Len(t, events, 3) {
		assert.True(t, events[0].Meta.EventDate.After(events[1].Meta.EventDate))
		assert.True(t, events[1].Meta.EventDate.After(events[2].Meta.EventDate))
	}

}

func TestGetEvents_ErrInvalidUserID(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	meta := &models.Meta{UserID: 0}

	_, err := service.GetEvents(meta, models.Day)
	assert.ErrorIs(t, err, errs.ErrInvalidUserID)

}

func TestGetEvents_ErrStorage(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)
	meta := &models.Meta{UserID: 1, EventDate: time.Now().UTC().Add(24 * time.Hour)}

	mockStorage.EXPECT().GetEvents(meta, models.Day).Return(nil, assert.AnError)

	events, err := service.GetEvents(meta, models.Day)
	assert.Nil(t, events)
	assert.ErrorIs(t, err, assert.AnError)

}

func TestGetEvents_ErrMissingDate(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := loggerMock.NewMockLogger(controller)
	mockStorage := storageMock.NewMockStorage(controller)

	service := NewService(config.Service{MaxEventsPerUser: 5}, mockStorage, mockLogger)

	meta := &models.Meta{
		UserID:    1,
		EventDate: time.Time{},
	}

	_, err := service.GetEvents(meta, models.Day)
	assert.ErrorIs(t, err, errs.ErrMissingDate)

}
