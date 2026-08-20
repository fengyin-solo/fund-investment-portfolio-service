package service

import (
	"sort"

	"fundinvest/internal/model"
	"fundinvest/pkg/idgen"
	"time"
)

func (s *Service) CreateHolding(input model.Holding) (*model.Holding, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetAccount(input.AccountID); err != nil {
		return nil, model.NewValidationError("account_id", "关联账户不存在")
	}
	if _, err := s.store.GetFund(input.FundID); err != nil {
		return nil, model.NewValidationError("fund_id", "关联基金不存在")
	}
	now := time.Now()
	h := &model.Holding{
		ID:          idgen.Hex(),
		AccountID:   input.AccountID,
		FundID:      input.FundID,
		TotalShares: input.TotalShares,
		Cost:        input.Cost,
		MarketValue: input.MarketValue,
		UpdatedAt:   now,
	}
	if err := s.store.CreateHolding(h); err != nil {
		return nil, err
	}
	return h, nil
}

func (s *Service) GetHolding(id string) (*model.Holding, error) {
	return s.store.GetHolding(id)
}

func (s *Service) ListHoldings(filter model.HoldingFilter, page, size int) ([]*model.Holding, int, error) {
	all := s.store.ListHoldings()
	matched := make([]*model.Holding, 0, len(all))
	for _, h := range all {
		if filter.Match(h) {
			matched = append(matched, h)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].UpdatedAt.After(matched[j].UpdatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Holding{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateHolding(id string, input model.Holding) (*model.Holding, error) {
	h, err := s.store.GetHolding(id)
	if err != nil {
		return nil, err
	}
	if input.TotalShares >= 0 {
		h.TotalShares = input.TotalShares
	}
	if input.Cost >= 0 {
		h.Cost = input.Cost
	}
	if input.MarketValue >= 0 {
		h.MarketValue = input.MarketValue
	}
	h.UpdatedAt = time.Now()
	if err := h.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateHolding(h); err != nil {
		return nil, err
	}
	return h, nil
}

func (s *Service) DeleteHolding(id string) error {
	return s.store.DeleteHolding(id)
}
