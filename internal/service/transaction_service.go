package service

import (
	"sort"
	"time"

	"fundinvest/internal/model"
	"fundinvest/pkg/idgen"
)

func (s *Service) CreateTransaction(input model.Transaction) (*model.Transaction, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetAccount(input.AccountID); err != nil {
		return nil, model.NewValidationError("account_id", "关联账户不存在")
	}
	if _, err := s.store.GetFund(input.FundID); err != nil {
		return nil, model.NewValidationError("fund_id", "关联基金不存在")
	}
	if input.PlanID != "" {
		if _, err := s.store.GetPlan(input.PlanID); err != nil {
			return nil, model.NewValidationError("plan_id", "关联计划不存在")
		}
	}
	t := &model.Transaction{
		ID:        idgen.Hex(),
		AccountID: input.AccountID,
		FundID:    input.FundID,
		PlanID:    input.PlanID,
		Type:      input.Type,
		Amount:    input.Amount,
		Shares:    input.Shares,
		Price:     input.Price,
		Status:    input.Status,
		CreatedAt: time.Now(),
	}
	if t.Status == "" {
		t.Status = model.TransactionStatusPending
	}
	if err := s.store.CreateTransaction(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetTransaction(id string) (*model.Transaction, error) {
	return s.store.GetTransaction(id)
}

func (s *Service) ListTransactions(filter model.TransactionFilter, page, size int) ([]*model.Transaction, int, error) {
	all := s.store.ListTransactions()
	matched := make([]*model.Transaction, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Transaction{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateTransaction(id string, input model.Transaction) (*model.Transaction, error) {
	t, err := s.store.GetTransaction(id)
	if err != nil {
		return nil, err
	}
	if input.Status != "" {
		t.Status = input.Status
	}
	if input.Amount >= 0 {
		t.Amount = input.Amount
	}
	if input.Shares >= 0 {
		t.Shares = input.Shares
	}
	if input.Price > 0 {
		t.Price = input.Price
	}
	if err := t.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateTransaction(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteTransaction(id string) error {
	return s.store.DeleteTransaction(id)
}

func (s *Service) Redeem(accountID, fundID string, shares int64) (*model.Transaction, error) {
	if shares <= 0 {
		return nil, model.NewValidationError("shares", "赎回份额必须大于 0")
	}
	holding, err := s.store.GetHoldingByAccountAndFund(accountID, fundID)
	if err != nil {
		return nil, model.NewValidationError("holding", "持仓不存在")
	}
	if holding.TotalShares < shares {
		return nil, model.NewValidationError("shares", "持仓份额不足")
	}
	fund, err := s.store.GetFund(fundID)
	if err != nil {
		return nil, err
	}
	if fund.Status != model.FundStatusActive {
		return nil, model.NewValidationError("fund_id", "基金未处于 active 状态")
	}
	account, err := s.store.GetAccount(accountID)
	if err != nil {
		return nil, err
	}

	amount := (shares * fund.Nav) / 10000
	now := time.Now()

	account.Balance += amount
	if err := s.store.UpdateAccount(account); err != nil {
		return nil, err
	}

	holding.TotalShares -= shares
	holding.Cost = (holding.Cost * holding.TotalShares) / (holding.TotalShares + shares)
	if holding.TotalShares == 0 {
		holding.Cost = 0
	}
	holding.MarketValue = (holding.TotalShares * fund.Nav) / 10000
	holding.UpdatedAt = now
	if err := s.store.UpdateHolding(holding); err != nil {
		return nil, err
	}

	txn := &model.Transaction{
		ID:        idgen.Hex(),
		AccountID: accountID,
		FundID:    fundID,
		Type:      model.TransactionTypeRedeem,
		Amount:    amount,
		Shares:    shares,
		Price:     fund.Nav,
		Status:    model.TransactionStatusConfirmed,
		CreatedAt: now,
	}
	if err := s.store.CreateTransaction(txn); err != nil {
		return nil, err
	}
	return txn, nil
}
