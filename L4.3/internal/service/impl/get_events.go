package impl

import (
	"sort"

	"L4.3/internal/models"
)

// GetEvents retrieves all events for a user within the specified period (day, week, month).
// Events are returned in descending order by date. Returns an error if validation fails
// or if the repository fails to fetch events.
func (s *Service) GetEvents(meta *models.Meta, period models.Period) ([]models.Event, error) {

	if err := validateGet(meta); err != nil {
		return nil, err
	}

	events, err := s.Storage.Memory.GetEvents(meta, period)
	if err != nil {
		return nil, err
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Meta.EventDate.After(events[j].Meta.EventDate)
	})

	return events, nil

}
