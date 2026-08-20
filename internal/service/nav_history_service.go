package service

import (
	"sort"
	"time"

	"fundinvest/internal/model"
	"fundinvest/pkg/idgen"
)

func (s *Service) CreateNavHistory(input model.NavHistory) (*model.NavHistory, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFund(input.FundID); err != nil {
		return nil, model.NewValidationError("fund_id", "关联基金不存在")
	}
	n := &model.NavHistory{
		ID:         idgen.Hex(),
		FundID:     input.FundID,
		Nav:        input.Nav,
		RecordedAt: time.Now(),
	}
	if !input.RecordedAt.IsZero() {
		n.RecordedAt = input.RecordedAt
	}
	if err := s.store.CreateNavHistory(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) GetNavHistory(id string) (*model.NavHistory, error) {
	return s.store.GetNavHistory(id)
}

func (s *Service) ListNavHistory(filter model.NavHistoryFilter, page, size int) ([]*model.NavHistory, int, error) {
	all := s.store.ListNavHistory()
	matched := make([]*model.NavHistory, 0, len(all))
	for _, n := range all {
		if filter.Match(n) {
			matched = append(matched, n)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].RecordedAt.After(matched[j].RecordedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.NavHistory{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateNavHistory(id string, input model.NavHistory) (*model.NavHistory, error) {
	n, err := s.store.GetNavHistory(id)
	if err != nil {
		return nil, err
	}
	if input.Nav > 0 {
		n.Nav = input.Nav
	}
	if !input.RecordedAt.IsZero() {
		n.RecordedAt = input.RecordedAt
	}
	if err := n.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateNavHistory(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) DeleteNavHistory(id string) error {
	return s.store.DeleteNavHistory(id)
}
