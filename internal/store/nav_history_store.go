package store

import "fundinvest/internal/model"

func (s *MemoryStore) CreateNavHistory(n *model.NavHistory) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.navHistory[n.ID] = n
	return nil
}

func (s *MemoryStore) GetNavHistory(id string) (*model.NavHistory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.navHistory[id]
	if !ok {
		return nil, ErrNotFound
	}
	return n, nil
}

func (s *MemoryStore) ListNavHistory() []*model.NavHistory {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.NavHistory, 0, len(s.navHistory))
	for _, n := range s.navHistory {
		list = append(list, n)
	}
	return list
}

func (s *MemoryStore) UpdateNavHistory(n *model.NavHistory) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.navHistory[n.ID]; !ok {
		return ErrNotFound
	}
	s.navHistory[n.ID] = n
	return nil
}

func (s *MemoryStore) DeleteNavHistory(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.navHistory[id]; !ok {
		return ErrNotFound
	}
	delete(s.navHistory, id)
	return nil
}
