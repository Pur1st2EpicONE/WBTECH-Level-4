package v1

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"L4.3/internal/errs"

	"github.com/stretchr/testify/assert"
)

func TestParseQuery_Success(t *testing.T) {

	userID := "8"
	dateStr := "2025-12-03"

	id, date, err := parseQuery(userID, dateStr)

	assert.NoError(t, err)
	assert.Equal(t, 8, id)

	expectedDate := time.Date(2025, 12, 3, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expectedDate, date)

}

func TestParseQuery_MissingParams(t *testing.T) {
	_, _, err := parseQuery("", "")
	assert.ErrorIs(t, err, errs.ErrMissingParams)
}

func TestParseQuery_InvalidUserID(t *testing.T) {
	_, _, err := parseQuery("qwe", "2025-12-03")
	assert.ErrorIs(t, err, errs.ErrInvalidUserID)
}

func TestParseQuery_InvalidDate(t *testing.T) {
	_, _, err := parseQuery("1", "2025-13-03")
	assert.ErrorIs(t, err, errs.ErrInvalidDateFormat)
}

func TestParseDate_Success(t *testing.T) {

	dateStr := "2025-12-03"
	date, err := parseDate(dateStr)

	assert.NoError(t, err)

	expectedDate := time.Date(2025, 12, 3, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expectedDate, date)

}

func TestParseDate_Empty(t *testing.T) {
	_, err := parseDate("")
	assert.ErrorIs(t, err, errs.ErrMissingDate)
}

func TestParseDate_InvalidFormat(t *testing.T) {
	_, err := parseDate("03-12-2025")
	assert.ErrorIs(t, err, errs.ErrInvalidDateFormat)
}

func TestBadRequest(t *testing.T) {

	tests := []error{
		errs.ErrInvalidJSON,
		errs.ErrInvalidUserID,
		errs.ErrInvalidEventID,
		errs.ErrInvalidDateFormat,
		errs.ErrEmptyEventText,
		errs.ErrEventTextTooLong,
		errs.ErrMissingEventID,
		errs.ErrMissingParams,
		errs.ErrMissingDate,
	}

	for _, e := range tests {
		status, msg := mapErrorToStatus(e)
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, e.Error(), msg)
	}

}

func TestServiceUnavailable(t *testing.T) {

	tests := []error{
		errs.ErrMaxEvents,
		errs.ErrEventNotFound,
		errs.ErrNothingToUpdate,
		errs.ErrEventInPast,
		errs.ErrEventTooFar,
		errs.ErrUnauthorized,
	}

	for _, e := range tests {
		status, msg := mapErrorToStatus(e)
		assert.Equal(t, http.StatusServiceUnavailable, status)
		assert.Equal(t, e.Error(), msg)
	}

}

func TestInternalServerError(t *testing.T) {
	err := errors.New("some bad error")
	status, msg := mapErrorToStatus(err)
	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, errs.ErrInternal.Error(), msg)
}
