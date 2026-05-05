package impl

import (
	"fmt"
	"time"

	"L4.3/internal/models"
)

// reminderWorker is a background goroutine responsible for processing reminders.
//
// It maintains an in-memory schedule of reminder jobs and periodically checks
// whether they are due.
func (s *Service) reminderWorker() {

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	jobs := make(map[string]models.Reminder)

	for {
		select {

		case job := <-s.reminderCh:
			jobs[job.EventID] = job

		case <-ticker.C:
			now := time.Now().UTC()

			for id, job := range jobs {
				if now.After(job.RemindAt) || now.Equal(job.RemindAt) {

					fmt.Printf("[REMINDER] user=%d event=%s text=%s\n",
						job.UserID, job.EventID, job.Text)

					delete(jobs, id)
				}
			}

		case <-s.stopCh:
			return
		}
	}
}
