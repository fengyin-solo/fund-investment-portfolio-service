package store

import "fundinvest/internal/model"

func (s *MemoryStore) CreateAccount(a *model.Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accounts[a.ID] = a
	return nil
}

func (s *MemoryStore) GetAccount(id string) (*model.Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.accounts[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) ListAccounts() []*model.Account {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Account, 0, len(s.accounts))
	for _, a := range s.accounts {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) UpdateAccount(a *model.Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.accounts[a.ID]; !ok {
		return ErrNotFound
	}
	s.accounts[a.ID] = a
	return nil
}

func (s *MemoryStore) DeleteAccount(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.accounts[id]; !ok {
		return ErrNotFound
	}
	delete(s.accounts, id)
	return nil
}
