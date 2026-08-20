package model

import (
	"time"
)

type NavHistory struct {
	ID         string    `json:"id"`
	FundID     string    `json:"fund_id"`
	Nav        int64     `json:"nav"`
	RecordedAt time.Time `json:"recorded_at"`
}

func (n *NavHistory) Validate() error {
	if n.FundID == "" {
		return NewValidationError("fund_id", "基金 ID 不能为空")
	}
	if n.Nav <= 0 {
		return NewValidationError("nav", "净值必须大于 0")
	}
	return nil
}

type NavHistoryFilter struct {
	FundID    string
	StartTime time.Time
	EndTime   time.Time
}

func (f NavHistoryFilter) Match(n *NavHistory) bool {
	if f.FundID != "" && n.FundID != f.FundID {
		return false
	}
	if !f.StartTime.IsZero() && n.RecordedAt.Before(f.StartTime) {
		return false
	}
	if !f.EndTime.IsZero() && n.RecordedAt.After(f.EndTime) {
		return false
	}
	return true
}
