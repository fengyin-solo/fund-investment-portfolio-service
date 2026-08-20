package store

import "fundinvest/internal/model"

func (s *MemoryStore) CreatePlan(p *model.Plan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.plans[p.ID] = p
	return nil
}

func (s *MemoryStore) GetPlan(id string) (*model.Plan, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.plans[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *MemoryStore) ListPlans() []*model.Plan {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Plan, 0, len(s.plans))
	for _, p := range s.plans {
		list = append(list, p)
	}
	return list
}

func (s *MemoryStore) UpdatePlan(p *model.Plan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.plans[p.ID]; !ok {
		return ErrNotFound
	}
	s.plans[p.ID] = p
	return nil
}

func (s *MemoryStore) DeletePlan(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.plans[id]; !ok {
		return ErrNotFound
	}
	delete(s.plans, id)
	return nil
}
