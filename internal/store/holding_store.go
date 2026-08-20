package store

import "fundinvest/internal/model"

func (s *MemoryStore) CreateHolding(h *model.Holding) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.holdings[h.ID] = h
	return nil
}

func (s *MemoryStore) GetHolding(id string) (*model.Holding, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h, ok := s.holdings[id]
	if !ok {
		return nil, ErrNotFound
	}
	return h, nil
}

func (s *MemoryStore) GetHoldingByAccountAndFund(accountID, fundID string) (*model.Holding, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, h := range s.holdings {
		if h.AccountID == accountID && h.FundID == fundID {
			return h, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListHoldings() []*model.Holding {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Holding, 0, len(s.holdings))
	for _, h := range s.holdings {
		list = append(list, h)
	}
	return list
}

func (s *MemoryStore) UpdateHolding(h *model.Holding) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.holdings[h.ID]; !ok {
		return ErrNotFound
	}
	s.holdings[h.ID] = h
	return nil
}

func (s *MemoryStore) DeleteHolding(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.holdings[id]; !ok {
		return ErrNotFound
	}
	delete(s.holdings, id)
	return nil
}
