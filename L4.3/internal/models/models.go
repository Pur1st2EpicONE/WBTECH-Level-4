// Package models defines the domain models for the application,
// including event representations and periods for filtering.
package models

import "time"

type Period string // Period represents a time period used for filtering events.

const (
	Day   Period = "day"   // Day represents a single day period.
	Week  Period = "week"  // Week represents a single week period.
	Month Period = "month" // Month represents a single month period.
)

// Event represents a user's event with metadata and associated data.
type Event struct {
	Meta Meta // Metadata about the event (ID, user, date)
	Data Data // Event-specific data (text)
}

// Meta contains identifying and timing information for an event.
type Meta struct {
	UserID    int       // ID of the user who owns the event
	EventID   string    // Unique identifier for the event
	EventDate time.Time // Original date of the event
	NewDate   time.Time // Updated date of the event (if modified)
}

// Data contains the actual content of the event.
type Data struct {
	Text string // Text description of the event
}
