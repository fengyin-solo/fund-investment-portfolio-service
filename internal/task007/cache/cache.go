package cache

import "fundinvest/internal/task007/model"

type Store struct{ jobs map[string]model.Job }

func New() *Store                         { return &Store{jobs: map[string]model.Job{}} }
func (s *Store) Apply(job model.Job) bool { s.jobs[job.ID] = job; return true }
func (s *Store) Get(id string) model.Job  { return s.jobs[id] }
