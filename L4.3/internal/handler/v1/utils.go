package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"L4.3/internal/errs"
	"github.com/gin-gonic/gin"
)

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

func parseDate(date string) (time.Time, error) {

	if date == "" {
		return time.Time{}, errs.ErrMissingDate
	}

	if t, err := time.Parse(time.RFC3339, date); err == nil {
		return t, nil
	}

	if t, err := time.Parse("2006-01-02", date); err == nil {
		return t, nil
	}

	return time.Time{}, errs.ErrInvalidDateFormat

}

func respondOK(c *gin.Context, response any) {
	c.JSON(http.StatusOK, gin.H{"result": response})
}

func respondError(c *gin.Context, err error) {
	if err != nil {
		status, msg := mapErrorToStatus(err)
		c.AbortWithStatusJSON(status, gin.H{"error": msg})
	}
}

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
