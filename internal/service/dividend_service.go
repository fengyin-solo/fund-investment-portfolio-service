package service

import (
	"sort"
	"time"

	"fundinvest/internal/model"
	"fundinvest/pkg/idgen"
)

func (s *Service) CreateDividend(input model.Dividend) (*model.Dividend, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetAccount(input.AccountID); err != nil {
		return nil, model.NewValidationError("account_id", "关联账户不存在")
	}
	if _, err := s.store.GetFund(input.FundID); err != nil {
		return nil, model.NewValidationError("fund_id", "关联基金不存在")
	}
	d := &model.Dividend{
		ID:        idgen.Hex(),
		AccountID: input.AccountID,
		FundID:    input.FundID,
		Amount:    input.Amount,
		Type:      input.Type,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateDividend(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) GetDividend(id string) (*model.Dividend, error) {
	return s.store.GetDividend(id)
}

func (s *Service) ListDividends(filter model.DividendFilter, page, size int) ([]*model.Dividend, int, error) {
	all := s.store.ListDividends()
	matched := make([]*model.Dividend, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Dividend{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateDividend(id string, input model.Dividend) (*model.Dividend, error) {
	d, err := s.store.GetDividend(id)
	if err != nil {
		return nil, err
	}
	if input.Amount > 0 {
		d.Amount = input.Amount
	}
	if input.Type != "" {
		d.Type = input.Type
	}
	if err := d.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateDividend(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) DeleteDividend(id string) error {
	return s.store.DeleteDividend(id)
}
