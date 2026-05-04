package impl

import (
	"errors"
	"testing"
	"time"

	"L4.3/internal/errs"
	"L4.3/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateDate_Past(t *testing.T) {
	past := time.Now().UTC().AddDate(0, 0, -1)
	err := validateDate(past)
	assert.ErrorIs(t, err, errs.ErrEventInPast)
}

func TestValidateDate_TooFar(t *testing.T) {
	tooFar := time.Now().UTC().AddDate(11, 0, 0)
	err := validateDate(tooFar)
	assert.True(t, errors.Is(err, errs.ErrEventTooFar))
}

func TestValidateDate_Success(t *testing.T) {
	ok := time.Now().UTC().AddDate(1, 0, 0)
	err := validateDate(ok)
	assert.NoError(t, err)
}

func TestValidateData_TooLong(t *testing.T) {
	long := models.Data{Text: string(make([]byte, 501))}
	err := validateData(long)
	assert.ErrorIs(t, err, errs.ErrEventTextTooLong)
}

func TestValidateData_Success(t *testing.T) {
	ok := models.Data{Text: "short"}
	err := validateData(ok)
	assert.NoError(t, err)
}

func TestValidateIDs_InvalidUserID(t *testing.T) {
	err := validateIDs(0, uuid.New().String())
	assert.ErrorIs(t, err, errs.ErrInvalidUserID)
}

func TestValidateIDs_MissingEventID(t *testing.T) {
	err := validateIDs(1, "")
	assert.ErrorIs(t, err, errs.ErrMissingEventID)
}

func TestValidateIDs_InvalidUUID(t *testing.T) {
	err := validateIDs(1, "not-a-uuid")
	assert.ErrorIs(t, err, errs.ErrInvalidEventID)
}

func TestValidateCreate_InvalidUserID(t *testing.T) {
	event := &models.Event{
		Meta: models.Meta{UserID: 0, EventDate: time.Now().Add(24 * time.Hour)},
		Data: models.Data{Text: "qwe"},
	}
	err := validateCreate(event)
	assert.ErrorIs(t, err, errs.ErrInvalidUserID)
}

func TestValidateCreate_InvalidDate(t *testing.T) {
	event := &models.Event{
		Meta: models.Meta{UserID: 1, EventDate: time.Now().AddDate(-1, 0, 0)},
		Data: models.Data{Text: "qwe"},
	}
	err := validateCreate(event)
	assert.True(t, errors.Is(err, errs.ErrEventInPast))
}

func TestValidateCreate_InvalidData(t *testing.T) {
	event := &models.Event{
		Meta: models.Meta{UserID: 1, EventDate: time.Now().Add(24 * time.Hour)},
		Data: models.Data{Text: string(make([]byte, 501))},
	}
	err := validateCreate(event)
	assert.ErrorIs(t, err, errs.ErrEventTextTooLong)
}

func TestValidateUpdate_ErrInvalidNewDate(t *testing.T) {

	now := time.Now().UTC().Add(24 * time.Hour)

	oldEvent := &models.Event{
		Meta: models.Meta{UserID: 1, EventDate: now},
		Data: models.Data{Text: "ok"},
	}

	event := &models.Event{
		Meta: models.Meta{UserID: 1, NewDate: time.Now().UTC().AddDate(-1, 0, 0)},
		Data: models.Data{Text: "ok"},
	}

	err := validateUpdate(event, oldEvent)
	assert.True(t, errors.Is(err, errs.ErrEventInPast))

}

func TestValidateUpdate_ErrInvalidData(t *testing.T) {

	now := time.Now().UTC().Add(24 * time.Hour)

	oldEvent := &models.Event{
		Meta: models.Meta{UserID: 1, EventDate: now},
		Data: models.Data{Text: "ok"},
	}

	event := &models.Event{
		Meta: models.Meta{UserID: 1, NewDate: now.Add(24 * time.Hour)},
		Data: models.Data{Text: string(make([]byte, 501))},
	}

	err := validateUpdate(event, oldEvent)
	assert.ErrorIs(t, err, errs.ErrEventTextTooLong)

}

func TestIsNothingToUpdate(t *testing.T) {

	now := time.Now().UTC()

	oldEvent := &models.Event{
		Meta: models.Meta{EventDate: now},
		Data: models.Data{Text: "same"},
	}

	t.Run("no changes", func(t *testing.T) {
		newEvent := &models.Event{
			Meta: models.Meta{EventDate: now},
			Data: models.Data{Text: "same"},
		}
		assert.True(t, isNothingToUpdate(newEvent, oldEvent))
	})

	t.Run("text changed", func(t *testing.T) {
		newEvent := &models.Event{
			Meta: models.Meta{EventDate: now},
			Data: models.Data{Text: "changed"},
		}
		assert.False(t, isNothingToUpdate(newEvent, oldEvent))
	})

	t.Run("date changed", func(t *testing.T) {
		newEvent := &models.Event{
			Meta: models.Meta{EventDate: oldEvent.Meta.EventDate, NewDate: oldEvent.Meta.EventDate.Add(24 * time.Hour)},
			Data: models.Data{Text: "same"},
		}
		assert.False(t, isNothingToUpdate(newEvent, oldEvent))
	})

}
