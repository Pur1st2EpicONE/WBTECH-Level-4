package impl

import "L4.3/internal/models"

// DeleteEvent validates and deletes an event identified by the provided metadata.
// Returns an error if validation fails or the event cannot be deleted.
func (s *Service) DeleteEvent(meta *models.Meta) error {

	if err := validateIDs(meta.UserID, meta.EventID); err != nil {
		return err
	}

	if err := validateDelete(meta, s.Storage.Memory.GetEventByID(meta.EventID)); err != nil {
		return err
	}

	return s.Storage.Memory.DeleteEvent(meta)

}
