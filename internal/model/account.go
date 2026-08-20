package model

import (
	"strings"
	"time"
)

type Account struct {
	ID        string    `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Account) Validate() error {
	a.Owner = strings.TrimSpace(a.Owner)
	if a.Owner == "" {
		return NewValidationError("owner", "账户持有人姓名不能为空")
	}
	if a.Balance < 0 {
		return NewValidationError("balance", "账户余额不能为负数")
	}
	return nil
}

type AccountFilter struct {
	Keyword string
}

func (f AccountFilter) Match(a *Account) bool {
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(a.Owner), k) {
			return false
		}
	}
	return true
}
