package model

import (
	"time"
)

type Holding struct {
	ID           string    `json:"id"`
	AccountID    string    `json:"account_id"`
	FundID       string    `json:"fund_id"`
	TotalShares  int64     `json:"total_shares"`
	Cost         int64     `json:"cost"`
	MarketValue  int64     `json:"market_value"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (h *Holding) Validate() error {
	if h.AccountID == "" {
		return NewValidationError("account_id", "账户 ID 不能为空")
	}
	if h.FundID == "" {
		return NewValidationError("fund_id", "基金 ID 不能为空")
	}
	if h.TotalShares < 0 {
		return NewValidationError("total_shares", "总份额不能为负数")
	}
	if h.Cost < 0 {
		return NewValidationError("cost", "成本不能为负数")
	}
	return nil
}

type HoldingFilter struct {
	AccountID string
	FundID    string
}

func (f HoldingFilter) Match(h *Holding) bool {
	if f.AccountID != "" && h.AccountID != f.AccountID {
		return false
	}
	if f.FundID != "" && h.FundID != f.FundID {
		return false
	}
	return true
}
