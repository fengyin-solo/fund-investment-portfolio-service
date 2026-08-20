package store

import "fundinvest/internal/model"

func (s *MemoryStore) CreateFund(f *model.Fund) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.funds {
		if exist.Code == f.Code {
			return ErrConflict
		}
	}
	s.funds[f.ID] = f
	return nil
}

func (s *MemoryStore) GetFund(id string) (*model.Fund, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.funds[id]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

func (s *MemoryStore) GetFundByCode(code string) (*model.Fund, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, f := range s.funds {
		if f.Code == code {
			return f, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListFunds() []*model.Fund {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Fund, 0, len(s.funds))
	for _, f := range s.funds {
		list = append(list, f)
	}
	return list
}

func (s *MemoryStore) UpdateFund(f *model.Fund) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.funds[f.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.funds {
		if exist.ID != f.ID && exist.Code == f.Code {
			return ErrConflict
		}
	}
	s.funds[f.ID] = f
	return nil
}

func (s *MemoryStore) DeleteFund(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.funds[id]; !ok {
		return ErrNotFound
	}
	delete(s.funds, id)
	return nil
}
