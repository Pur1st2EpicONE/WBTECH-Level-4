package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"L4.3/internal/errs"
	"github.com/gin-gonic/gin"
)

// parseQuery parses and validates the query parameters from an HTTP request.
//
// userID: string representing the user's ID from query parameters.
// eventDate: string representing the event date in YYYY-MM-DD format.
//
// Returns:
// - user ID as int
// - parsed event date as time.Time
// - error if any validation fails (missing params, invalid ID, invalid date format)
func parseQuery(userID string, eventDate string) (int, time.Time, error) {

	if userID == "" || eventDate == "" {
		return 0, time.Time{}, errs.ErrMissingParams
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, time.Time{}, errs.ErrInvalidUserID
	}

	date, err := parseDate(eventDate)
	if err != nil {
		return 0, time.Time{}, err
	}

	return id, date, nil

}

// parseDate parses a date string in "YYYY-MM-DD" format.
//
// date: string representation of the date.
//
// Returns:
// - parsed date as time.Time
// - error if the date is missing or the format is invalid
func parseDate(date string) (time.Time, error) {

	if date == "" {
		return time.Time{}, errs.ErrMissingDate
	}

	eventDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, errs.ErrInvalidDateFormat
	}

	return eventDate, nil

}

// respondOK sends a successful JSON response to the client.
//
// c: Gin context
// response: the response payload to send
func respondOK(c *gin.Context, response any) {
	c.JSON(http.StatusOK, gin.H{"result": response})
}

// respondError sends an error JSON response to the client based on the error type.
//
// c: Gin context
// err: the error to map and send
func respondError(c *gin.Context, err error) {
	if err != nil {
		status, msg := mapErrorToStatus(err)
		c.AbortWithStatusJSON(status, gin.H{"error": msg})
	}
}

// mapErrorToStatus maps application errors to HTTP status codes and messages.
//
// err: the error to map
//
// Returns:
// - HTTP status code
// - error message string to send in response
func mapErrorToStatus(err error) (int, string) {

	switch {

	case errors.Is(err, errs.ErrInvalidJSON),
		errors.Is(err, errs.ErrInvalidUserID),
		errors.Is(err, errs.ErrInvalidEventID),
		errors.Is(err, errs.ErrInvalidDateFormat),
		errors.Is(err, errs.ErrEmptyEventText),
		errors.Is(err, errs.ErrEventTextTooLong),
		errors.Is(err, errs.ErrMissingEventID),
		errors.Is(err, errs.ErrMissingParams),
		errors.Is(err, errs.ErrMissingDate):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, errs.ErrMaxEvents),
		errors.Is(err, errs.ErrEventNotFound),
		errors.Is(err, errs.ErrNothingToUpdate),
		errors.Is(err, errs.ErrEventInPast),
		errors.Is(err, errs.ErrEventTooFar),
		errors.Is(err, errs.ErrUnauthorized):
		return http.StatusServiceUnavailable, err.Error()

	default:
		return http.StatusInternalServerError, errs.ErrInternal.Error()

	}

}
