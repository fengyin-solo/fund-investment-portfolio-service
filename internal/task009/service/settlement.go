package service

import (
	"fmt"
	"fundinvest/internal/task009/audit"
	"fundinvest/internal/task009/publisher"
	"fundinvest/internal/task009/repository"
)

type Settlement struct {
	Repo  *repository.Repository
	Bus   *publisher.Bus
	Audit *audit.Log
}

func (s *Settlement) Run(key string) error {
	s.Audit.Begin(key)

	// The settlement write is durable and must happen exactly once: commit it
	// before publishing, so a publish failure never re-opens the transaction.
	tx := s.Repo.Begin(key)
	tx.Commit()

	// Only the unreliable event publish is retried. A stable, attempt-
	// independent key makes retries idempotent (the same event, not a new one).
	eventKey := fmt.Sprintf("settled:%s", key)
	for attempt := 0; attempt < 2; attempt++ {
		if err := s.Bus.Publish(eventKey); err != nil {
			if attempt == 1 {
				return err
			}
			continue
		}
		break
	}

	s.Audit.Complete(key)
	return nil
}
