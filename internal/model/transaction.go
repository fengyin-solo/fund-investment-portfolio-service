package model

import (
	"strings"
	"time"
)

const (
	TransactionTypeInvest = "invest"
	TransactionTypeRedeem = "redeem"
)

const (
	TransactionStatusPending   = "pending"
	TransactionStatusConfirmed = "confirmed"
	TransactionStatusFailed    = "failed"
)

type Transaction struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	FundID    string    `json:"fund_id"`
	PlanID    string    `json:"plan_id,omitempty"`
	Type      string    `json:"type"`
	Amount    int64     `json:"amount"`
	Shares    int64     `json:"shares"`
	Price     int64     `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Transaction) Validate() error {
	t.AccountID = strings.TrimSpace(t.AccountID)
	t.FundID = strings.TrimSpace(t.FundID)
	t.Type = strings.TrimSpace(t.Type)
	if t.AccountID == "" {
		return NewValidationError("account_id", "账户 ID 不能为空")
	}
	if t.FundID == "" {
		return NewValidationError("fund_id", "基金 ID 不能为空")
	}
	if t.Type == "" {
		return NewValidationError("type", "交易类型不能为空")
	}
	if t.Type != TransactionTypeInvest && t.Type != TransactionTypeRedeem {
		return NewValidationError("type", "交易类型不合法")
	}
	if t.Amount < 0 {
		return NewValidationError("amount", "交易金额不能为负数")
	}
	if t.Shares < 0 {
		return NewValidationError("shares", "份额不能为负数")
	}
	if t.Price <= 0 {
		return NewValidationError("price", "交易单价必须大于 0")
	}
	if t.Status == "" {
		t.Status = TransactionStatusPending
	}
	if t.Status != TransactionStatusPending && t.Status != TransactionStatusConfirmed && t.Status != TransactionStatusFailed {
		return NewValidationError("status", "交易状态不合法")
	}
	return nil
}

type TransactionFilter struct {
	AccountID  string
	FundID     string
	Type       string
	Status     string
	PlanID     string
	StartTime  time.Time
	EndTime    time.Time
}

func (f TransactionFilter) Match(t *Transaction) bool {
	if f.AccountID != "" && t.AccountID != f.AccountID {
		return false
	}
	if f.FundID != "" && t.FundID != f.FundID {
		return false
	}
	if f.Type != "" && t.Type != f.Type {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.PlanID != "" && t.PlanID != f.PlanID {
		return false
	}
	if !f.StartTime.IsZero() && t.CreatedAt.Before(f.StartTime) {
		return false
	}
	if !f.EndTime.IsZero() && t.CreatedAt.After(f.EndTime) {
		return false
	}
	return true
}
