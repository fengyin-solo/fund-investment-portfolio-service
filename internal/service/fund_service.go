package service

import (
	"sort"
	"time"

	"fundinvest/internal/model"
	"fundinvest/pkg/idgen"
)

func (s *Service) CreateFund(input model.Fund) (*model.Fund, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	f := &model.Fund{
		ID:        idgen.Hex(),
		Code:      input.Code,
		Name:      input.Name,
		Type:      input.Type,
		RiskLevel: input.RiskLevel,
		Nav:       input.Nav,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateFund(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Service) GetFund(id string) (*model.Fund, error) {
	return s.store.GetFund(id)
}

func (s *Service) ListFunds(filter model.FundFilter, page, size int) ([]*model.Fund, int, error) {
	all := s.store.ListFunds()
	matched := make([]*model.Fund, 0, len(all))
	for _, f := range all {
		if filter.Match(f) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Fund{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateFund(id string, input model.Fund) (*model.Fund, error) {
	f, err := s.store.GetFund(id)
	if err != nil {
		return nil, err
	}
	if input.Code != "" {
		f.Code = input.Code
	}
	if input.Name != "" {
		f.Name = input.Name
	}
	if input.Type != "" {
		f.Type = input.Type
	}
	if input.RiskLevel > 0 {
		f.RiskLevel = input.RiskLevel
	}
	if input.Nav > 0 {
		f.Nav = input.Nav
	}
	if input.Status != "" {
		f.Status = input.Status
	}
	f.UpdatedAt = time.Now()
	if err := f.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateFund(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Service) DeleteFund(id string) error {
	return s.store.DeleteFund(id)
}
