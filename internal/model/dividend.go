package model

import (
	"strings"
	"time"
)

const (
	DividendTypeCash    = "cash"
	DividendTypeReinvest = "reinvest"
)

type Dividend struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	FundID    string    `json:"fund_id"`
	Amount    int64     `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

func (d *Dividend) Validate() error {
	d.AccountID = strings.TrimSpace(d.AccountID)
	d.FundID = strings.TrimSpace(d.FundID)
	d.Type = strings.TrimSpace(d.Type)
	if d.AccountID == "" {
		return NewValidationError("account_id", "账户 ID 不能为空")
	}
	if d.FundID == "" {
		return NewValidationError("fund_id", "基金 ID 不能为空")
	}
	if d.Amount <= 0 {
		return NewValidationError("amount", "分红金额必须大于 0")
	}
	if d.Type == "" {
		return NewValidationError("type", "分红类型不能为空")
	}
	if d.Type != DividendTypeCash && d.Type != DividendTypeReinvest {
		return NewValidationError("type", "分红类型不合法")
	}
	return nil
}

type DividendFilter struct {
	AccountID string
	FundID    string
	Type      string
}

func (f DividendFilter) Match(d *Dividend) bool {
	if f.AccountID != "" && d.AccountID != f.AccountID {
		return false
	}
	if f.FundID != "" && d.FundID != f.FundID {
		return false
	}
	if f.Type != "" && d.Type != f.Type {
		return false
	}
	return true
}
