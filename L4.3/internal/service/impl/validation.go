package impl

import (
	"fmt"
	"time"

	"L4.3/internal/errs"
	"L4.3/internal/models"
	"github.com/google/uuid"
)

// validateCreate performs validation on a new event before creation.
// It checks that the user ID is valid, the event date is acceptable,
// and the event data meets length constraints.
func validateCreate(event *models.Event) error {

	if event.Meta.UserID <= 0 {
		return errs.ErrInvalidUserID
	}

	// if event.Data.Reminder < 0 {
	// 	event.Data.Reminder = 0
	// }

	if err := validateDate(event.Meta.EventDate); err != nil {
		return err
	}

	if err := validateData(event.Data); err != nil {
		return err
	}

	return nil

}

// validateUpdate checks whether an update to an existing event is valid.
// It ensures the event exists, belongs to the user, and that the update actually
// changes either the date or the text. It also validates any new date or text.
func validateUpdate(event *models.Event, oldEvent *models.Event) error {

	if oldEvent == nil {
		return errs.ErrEventNotFound
	}

	if event.Meta.UserID != oldEvent.Meta.UserID {
		return errs.ErrUnauthorized
	}

	if isNothingToUpdate(event, oldEvent) {
		return errs.ErrNothingToUpdate
	}

	if !event.Meta.NewDate.IsZero() {
		if err := validateDate(event.Meta.NewDate); err != nil {
			return err
		}
	}

	if event.Data.Text != "" {
		if err := validateData(event.Data); err != nil {
			return err
		}
	}

	return nil
}

// isNothingToUpdate returns true if the new event has neither a changed date
// nor changed text compared to the existing event.
func isNothingToUpdate(event *models.Event, oldEvent *models.Event) bool {
	if !event.Meta.NewDate.IsZero() && !oldEvent.Meta.EventDate.Equal(event.Meta.NewDate) {
		return false
	}
	if event.Data.Text != oldEvent.Data.Text {
		return false
	}
	return true
}

// validateDelete checks whether an event can be deleted.
// Returns an error if the event does not exist or does not belong to the user.
func validateDelete(meta *models.Meta, oldEvent *models.Event) error {

	if oldEvent == nil {
		return errs.ErrEventNotFound
	}

	if meta.UserID != oldEvent.Meta.UserID {
		return errs.ErrUnauthorized
	}

	return nil

}

// validateGet checks if the request to retrieve events is valid.
// UserID must be positive and EventDate must be set.
func validateGet(meta *models.Meta) error {

	if meta.UserID <= 0 {
		return errs.ErrInvalidUserID
	}

	if meta.EventDate.IsZero() {
		return errs.ErrMissingDate
	}

	return nil

}

// validateDate ensures the event date is not in the past and not more than 10 years ahead.
// Dates are compared in UTC to prevent timezone-related errors.
func validateDate(date time.Time) error {

	eventUTC := date.UTC().Truncate(24 * time.Hour)
	todayUTC := time.Now().UTC().Truncate(24 * time.Hour)

	if eventUTC.Before(todayUTC) {
		return fmt.Errorf("%w: %s", errs.ErrEventInPast, eventUTC.Format("2006-01-02"))

	}

	if eventUTC.After(todayUTC.AddDate(10, 0, 0)) {
		return fmt.Errorf("%w: %s", errs.ErrEventTooFar, eventUTC.Format("2006-01-02"))
	}

	return nil

}

// validateData encapsulates the validation logic for the event's data.
func validateData(data models.Data) error {

	if len(data.Text) > 500 {
		return errs.ErrEventTextTooLong
	}

	return nil

}

// validateIDs validates the user ID and event ID for update or delete operations.
// Checks include: positive userID, non-empty eventID, and proper UUID format for eventID.
func validateIDs(userID int, eventID string) error {

	if userID <= 0 {
		return errs.ErrInvalidUserID
	}

	if eventID == "" {
		return errs.ErrMissingEventID
	}

	_, err := uuid.Parse(eventID)
	if err != nil {
		return errs.ErrInvalidEventID
	}

	return nil

}
