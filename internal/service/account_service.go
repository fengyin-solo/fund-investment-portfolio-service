package service

import (
	"sort"
	"time"

	"fundinvest/internal/model"
	"fundinvest/pkg/idgen"
)

func (s *Service) CreateAccount(input model.Account) (*model.Account, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	a := &model.Account{
		ID:        idgen.Hex(),
		Owner:     input.Owner,
		Balance:   input.Balance,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateAccount(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) GetAccount(id string) (*model.Account, error) {
	return s.store.GetAccount(id)
}

func (s *Service) ListAccounts(filter model.AccountFilter, page, size int) ([]*model.Account, int, error) {
	all := s.store.ListAccounts()
	matched := make([]*model.Account, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Account{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateAccount(id string, input model.Account) (*model.Account, error) {
	a, err := s.store.GetAccount(id)
	if err != nil {
		return nil, err
	}
	if input.Owner != "" {
		a.Owner = input.Owner
	}
	if input.Balance >= 0 {
		a.Balance = input.Balance
	}
	if err := a.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateAccount(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) DeleteAccount(id string) error {
	return s.store.DeleteAccount(id)
}
