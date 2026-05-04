// Package errs defines application-specific error variables used across
// the event management system.
//
// Each error represents a specific failure case in handling requests,
// such as invalid input, missing parameters, unauthorized access, or
// internal server errors. These errors can be returned by handlers,
// services, or repositories to provide consistent error reporting.
package errs

import "errors"

var (
	ErrInvalidJSON       = errors.New("invalid JSON format")                                 // invalid JSON format
	ErrEmptyEventText    = errors.New("event text cannot be empty")                          // event text cannot be empty
	ErrMissingDate       = errors.New("event date is required")                              // event date is required
	ErrInvalidDateFormat = errors.New("invalid date format, expected YYYY-MM-DD")            // invalid date format, expected YYYY-MM-DD
	ErrEventTextTooLong  = errors.New("event text exceeds maximum length of 500 characters") // event text exceeds maximum length of 500 characters
	ErrInvalidUserID     = errors.New("missing or invalid user ID")                          // missing or invalid user ID
	ErrEventInPast       = errors.New("event date cannot be in the past")                    // event date cannot be in the past
	ErrEventTooFar       = errors.New("event date cannot be more than 10 years ahead")       // event date cannot be more than 10 years ahead
	ErrMaxEvents         = errors.New("maximum number of events reached")                    // maximum number of events reached
	ErrNothingToUpdate   = errors.New("no changes detected to update")                       // no changes detected to update
	ErrEventNotFound     = errors.New("event not found")                                     // event not found
	ErrInvalidEventID    = errors.New("invalid event ID format")                             // invalid event ID format
	ErrUnauthorized      = errors.New("unauthorized: you cannot modify this event")          // unauthorized: you cannot modify this event
	ErrMissingParams     = errors.New("missing required parameters: user_id or date")        // missing required parameters: user_id or date
	ErrMissingEventID    = errors.New("event ID is required")                                // event ID is required
	ErrInternal          = errors.New("internal server error")                               // internal server error
)
