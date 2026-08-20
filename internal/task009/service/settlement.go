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
	for attempt := 0; attempt < 2; attempt++ {
		tx := s.Repo.Begin(key)
		eventKey := fmt.Sprintf("settled:%s:%d", key, attempt)
		if err := s.Bus.Publish(eventKey); err != nil {
			tx.Rollback()
			if attempt == 1 {
				return err
			}
			continue
		}
		tx.Commit()
		s.Audit.Complete(key)
		return nil
	}
	return nil
}
