package model

import (
	"strings"
	"time"
)

const (
	FundTypeStock  = "stock"
	FundTypeBond   = "bond"
	FundTypeMixed  = "mixed"
	FundTypeIndex  = "index"
)

const (
	FundStatusActive   = "active"
	FundStatusInactive = "inactive"
)

type Fund struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	RiskLevel int       `json:"risk_level"`
	Nav       int64     `json:"nav"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *Fund) Validate() error {
	f.Code = strings.TrimSpace(f.Code)
	f.Name = strings.TrimSpace(f.Name)
	f.Type = strings.TrimSpace(f.Type)
	if f.Code == "" {
		return NewValidationError("code", "基金代码不能为空")
	}
	if f.Name == "" {
		return NewValidationError("name", "基金名称不能为空")
	}
	if f.Type == "" {
		return NewValidationError("type", "基金类型不能为空")
	}
	if f.Type != FundTypeStock && f.Type != FundTypeBond && f.Type != FundTypeMixed && f.Type != FundTypeIndex {
		return NewValidationError("type", "基金类型不合法")
	}
	if f.RiskLevel < 1 || f.RiskLevel > 5 {
		return NewValidationError("risk_level", "风险等级必须在 1-5 之间")
	}
	if f.Nav <= 0 {
		return NewValidationError("nav", "净值必须大于 0")
	}
	if f.Status == "" {
		f.Status = FundStatusActive
	}
	if f.Status != FundStatusActive && f.Status != FundStatusInactive {
		return NewValidationError("status", "基金状态不合法")
	}
	return nil
}

type FundFilter struct {
	Type     string
	Status   string
	RiskLevel int
	Keyword  string
}

func (f FundFilter) Match(fund *Fund) bool {
	if f.Type != "" && fund.Type != f.Type {
		return false
	}
	if f.Status != "" && fund.Status != f.Status {
		return false
	}
	if f.RiskLevel > 0 && fund.RiskLevel != f.RiskLevel {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(fund.Name), k) && !strings.Contains(strings.ToLower(fund.Code), k) {
			return false
		}
	}
	return true
}
