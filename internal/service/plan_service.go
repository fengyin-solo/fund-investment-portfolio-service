package service

import (
	"sort"
	"time"

	"fundinvest/internal/model"
	"fundinvest/pkg/idgen"
)

func (s *Service) CreatePlan(input model.Plan) (*model.Plan, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetAccount(input.AccountID); err != nil {
		return nil, model.NewValidationError("account_id", "关联账户不存在")
	}
	fund, err := s.store.GetFund(input.FundID)
	if err != nil {
		return nil, model.NewValidationError("fund_id", "关联基金不存在")
	}
	if fund.Status != model.FundStatusActive {
		return nil, model.NewValidationError("fund_id", "基金未处于 active 状态")
	}
	now := time.Now()
	p := &model.Plan{
		ID:            idgen.Hex(),
		AccountID:     input.AccountID,
		FundID:        input.FundID,
		Amount:        input.Amount,
		Interval:      input.Interval,
		Status:        input.Status,
		NextExecAt:    input.NextExecAt,
		InvestedCount: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if p.Status == "" {
		p.Status = model.PlanStatusActive
	}
	if p.NextExecAt.IsZero() {
		p.NextExecAt = now
	}
	if err := s.store.CreatePlan(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetPlan(id string) (*model.Plan, error) {
	return s.store.GetPlan(id)
}

func (s *Service) ListPlans(filter model.PlanFilter, page, size int) ([]*model.Plan, int, error) {
	all := s.store.ListPlans()
	matched := make([]*model.Plan, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Plan{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdatePlan(id string, input model.Plan) (*model.Plan, error) {
	p, err := s.store.GetPlan(id)
	if err != nil {
		return nil, err
	}
	if input.Amount > 0 {
		p.Amount = input.Amount
	}
	if input.Interval != "" {
		p.Interval = input.Interval
	}
	if !input.NextExecAt.IsZero() {
		p.NextExecAt = input.NextExecAt
	}
	p.UpdatedAt = time.Now()
	if err := p.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdatePlan(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ChangePlanStatus(id string, toStatus string) (*model.Plan, error) {
	p, err := s.store.GetPlan(id)
	if err != nil {
		return nil, err
	}
	if !model.PlanCanTransition(p.Status, toStatus) {
		return nil, model.NewValidationError("status", "状态转换不合法")
	}
	p.Status = toStatus
	p.UpdatedAt = time.Now()
	if err := s.store.UpdatePlan(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) DeletePlan(id string) error {
	return s.store.DeletePlan(id)
}

func (s *Service) ExecutePlan(planID string) (*model.Transaction, error) {
	p, err := s.store.GetPlan(planID)
	if err != nil {
		return nil, err
	}
	if p.Status != model.PlanStatusActive {
		return nil, model.NewValidationError("status", "计划非 active 状态")
	}
	now := time.Now()
	if p.NextExecAt.After(now) {
		return nil, model.NewValidationError("next_exec_at", "计划尚未到期")
	}
	account, err := s.store.GetAccount(p.AccountID)
	if err != nil {
		return nil, err
	}
	fund, err := s.store.GetFund(p.FundID)
	if err != nil {
		return nil, err
	}
	if fund.Status != model.FundStatusActive {
		return nil, model.NewValidationError("fund_id", "基金未处于 active 状态")
	}

	if account.Balance < p.Amount {
		txn := &model.Transaction{
			ID:        idgen.Hex(),
			AccountID: p.AccountID,
			FundID:    p.FundID,
			PlanID:    p.ID,
			Type:      model.TransactionTypeInvest,
			Amount:    p.Amount,
			Shares:    0,
			Price:     fund.Nav,
			Status:    model.TransactionStatusFailed,
			CreatedAt: now,
		}
		_ = s.store.CreateTransaction(txn)
		return txn, nil
	}

	account.Balance -= p.Amount
	if err := s.store.UpdateAccount(account); err != nil {
		return nil, err
	}

	shares := (p.Amount * 10000) / fund.Nav
	txn := &model.Transaction{
		ID:        idgen.Hex(),
		AccountID: p.AccountID,
		FundID:    p.FundID,
		PlanID:    p.ID,
		Type:      model.TransactionTypeInvest,
		Amount:    p.Amount,
		Shares:    shares,
		Price:     fund.Nav,
		Status:    model.TransactionStatusConfirmed,
		CreatedAt: now,
	}
	if err := s.store.CreateTransaction(txn); err != nil {
		return nil, err
	}

	holding, err := s.store.GetHoldingByAccountAndFund(p.AccountID, p.FundID)
	if err != nil {
		holding = &model.Holding{
			ID:          idgen.Hex(),
			AccountID:   p.AccountID,
			FundID:      p.FundID,
			TotalShares: 0,
			Cost:        0,
			MarketValue: 0,
			UpdatedAt:   now,
		}
		if err := s.store.CreateHolding(holding); err != nil {
			return nil, err
		}
	}
	holding.TotalShares += shares
	holding.Cost += p.Amount
	holding.MarketValue = (holding.TotalShares * fund.Nav) / 10000
	holding.UpdatedAt = now
	if err := s.store.UpdateHolding(holding); err != nil {
		return nil, err
	}

	p.InvestedCount++
	switch p.Interval {
	case model.PlanIntervalDaily:
		p.NextExecAt = p.NextExecAt.AddDate(0, 0, 1)
	case model.PlanIntervalWeekly:
		p.NextExecAt = p.NextExecAt.AddDate(0, 0, 7)
	case model.PlanIntervalMonthly:
		p.NextExecAt = p.NextExecAt.AddDate(0, 1, 0)
	}
	p.UpdatedAt = now
	if err := s.store.UpdatePlan(p); err != nil {
		return nil, err
	}

	return txn, nil
}
