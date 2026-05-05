package postgres

import (
	"fmt"
	"strings"

	"L4.3/internal/models"
)

func (s *Storage) SaveEvents(events []models.Event) error {

	if len(events) == 0 {
		return nil
	}

	tx, err := s.db.Master.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var queryBuilder strings.Builder
	var args []any

	queryBuilder.WriteString(`

	INSERT INTO archived_events (event_id, user_id, event_date, text)
	VALUES 

	`)
	for i, e := range events {

		offset := i * 4

		fmt.Fprintf(&queryBuilder, "($%d,$%d,$%d,$%d)", offset+1, offset+2, offset+3, offset+4)

		if i != len(events)-1 {
			queryBuilder.WriteString(",")
		}

		args = append(args,
			e.Meta.EventID,
			e.Meta.UserID,
			e.Meta.EventDate,
			e.Data.Text,
		)
	}

	queryBuilder.WriteString(" ON CONFLICT (event_id) DO NOTHING")

	_, err = tx.Exec(queryBuilder.String(), args...)
	if err != nil {
		return fmt.Errorf("insert archived events: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	s.logger.Debug("events saved to archive", "count", len(events), "layer", "repository.postgres")

	return nil

}
