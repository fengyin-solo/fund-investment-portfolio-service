package model

import (
	"strings"
	"time"
)

const (
	PlanIntervalDaily   = "daily"
	PlanIntervalWeekly  = "weekly"
	PlanIntervalMonthly = "monthly"
)

const (
	PlanStatusActive  = "active"
	PlanStatusPaused  = "paused"
	PlanStatusClosed  = "closed"
)

type Plan struct {
	ID            string    `json:"id"`
	AccountID     string    `json:"account_id"`
	FundID        string    `json:"fund_id"`
	Amount        int64     `json:"amount"`
	Interval      string    `json:"interval"`
	Status        string    `json:"status"`
	NextExecAt    time.Time `json:"next_exec_at"`
	InvestedCount int       `json:"invested_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (p *Plan) Validate() error {
	p.AccountID = strings.TrimSpace(p.AccountID)
	p.FundID = strings.TrimSpace(p.FundID)
	p.Interval = strings.TrimSpace(p.Interval)
	if p.AccountID == "" {
		return NewValidationError("account_id", "账户 ID 不能为空")
	}
	if p.FundID == "" {
		return NewValidationError("fund_id", "基金 ID 不能为空")
	}
	if p.Amount <= 0 {
		return NewValidationError("amount", "每期金额必须大于 0")
	}
	if p.Interval == "" {
		return NewValidationError("interval", "定投周期不能为空")
	}
	if p.Interval != PlanIntervalDaily && p.Interval != PlanIntervalWeekly && p.Interval != PlanIntervalMonthly {
		return NewValidationError("interval", "定投周期不合法")
	}
	if p.Status == "" {
		p.Status = PlanStatusActive
	}
	if p.Status != PlanStatusActive && p.Status != PlanStatusPaused && p.Status != PlanStatusClosed {
		return NewValidationError("status", "计划状态不合法")
	}
	return nil
}

var planTransitions = map[string]map[string]bool{
	PlanStatusActive: {PlanStatusPaused: true, PlanStatusClosed: true},
	PlanStatusPaused: {PlanStatusActive: true, PlanStatusClosed: true},
}

func PlanCanTransition(from, to string) bool {
	if m, ok := planTransitions[from]; ok {
		return m[to]
	}
	return false
}

type PlanFilter struct {
	AccountID string
	FundID    string
	Status    string
	Interval  string
}

func (f PlanFilter) Match(p *Plan) bool {
	if f.AccountID != "" && p.AccountID != f.AccountID {
		return false
	}
	if f.FundID != "" && p.FundID != f.FundID {
		return false
	}
	if f.Status != "" && p.Status != f.Status {
		return false
	}
	if f.Interval != "" && p.Interval != f.Interval {
		return false
	}
	return true
}
