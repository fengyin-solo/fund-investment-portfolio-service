package cache

import "fundinvest/internal/task007/model"

type Store struct{ jobs map[string]model.Job }

func New() *Store                         { return &Store{jobs: map[string]model.Job{}} }
func (s *Store) Apply(job model.Job) bool {
	// Stale delayed callbacks (e.g. the pre-retry snapshot) must not roll a
	// newer state back, so only accept a strictly newer version.
	if job.Version > s.jobs[job.ID].Version {
		s.jobs[job.ID] = job
		return true
	}
	return false
}
func (s *Store) Get(id string) model.Job  { return s.jobs[id] }
