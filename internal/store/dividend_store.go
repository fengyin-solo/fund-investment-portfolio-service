package store

import "fundinvest/internal/model"

func (s *MemoryStore) CreateDividend(d *model.Dividend) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dividends[d.ID] = d
	return nil
}

func (s *MemoryStore) GetDividend(id string) (*model.Dividend, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.dividends[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) ListDividends() []*model.Dividend {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Dividend, 0, len(s.dividends))
	for _, d := range s.dividends {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) UpdateDividend(d *model.Dividend) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dividends[d.ID]; !ok {
		return ErrNotFound
	}
	s.dividends[d.ID] = d
	return nil
}

func (s *MemoryStore) DeleteDividend(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dividends[id]; !ok {
		return ErrNotFound
	}
	delete(s.dividends, id)
	return nil
}
