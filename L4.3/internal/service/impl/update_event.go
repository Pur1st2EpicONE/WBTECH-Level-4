package impl

import "L4.3/internal/models"

// UpdateEvent validates and updates an existing event.
// Returns an error if validation fails, the event does not exist, or the update cannot be applied.
func (s *Service) UpdateEvent(event *models.Event) error {

	if err := validateIDs(event.Meta.UserID, event.Meta.EventID); err != nil {
		return err
	}

	if err := validateUpdate(event, s.Storage.Memory.GetEventByID(event.Meta.EventID)); err != nil {
		return err
	}

	return s.Storage.Memory.UpdateEvent(event)

}
