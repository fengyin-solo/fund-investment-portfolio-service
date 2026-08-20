package store

import "fundinvest/internal/model"

func (s *MemoryStore) CreateTransaction(t *model.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTransaction(id string) (*model.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.transactions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListTransactions() []*model.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Transaction, 0, len(s.transactions))
	for _, t := range s.transactions {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTransaction(t *model.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.transactions[t.ID]; !ok {
		return ErrNotFound
	}
	s.transactions[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTransaction(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.transactions[id]; !ok {
		return ErrNotFound
	}
	delete(s.transactions, id)
	return nil
}
