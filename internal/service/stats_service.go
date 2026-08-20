package service

import (
	"sort"

	"fundinvest/internal/model"
)

type HoldingSummary struct {
	FundID      string `json:"fund_id"`
	FundName    string `json:"fund_name"`
	FundCode    string `json:"fund_code"`
	TotalShares int64  `json:"total_shares"`
	Cost        int64  `json:"cost"`
	MarketValue int64  `json:"market_value"`
	Profit      int64  `json:"profit"`
}

type AccountProfit struct {
	AccountID    string `json:"account_id"`
	Owner        string `json:"owner"`
	TotalCost    int64  `json:"total_cost"`
	TotalMarket  int64  `json:"total_market"`
	TotalProfit  int64  `json:"total_profit"`
}

type FundTypeStat struct {
	Type      string `json:"type"`
	Count     int    `json:"count"`
	TotalNav  int64  `json:"total_nav"`
}

type FundRiskStat struct {
	RiskLevel int `json:"risk_level"`
	Count     int `json:"count"`
}

type TopFund struct {
	FundID        string `json:"fund_id"`
	FundName      string `json:"fund_name"`
	FundCode      string `json:"fund_code"`
	TotalInvested int64  `json:"total_invested"`
}

func (s *Service) GetAccountHoldingSummary(accountID string) ([]HoldingSummary, error) {
	if _, err := s.store.GetAccount(accountID); err != nil {
		return nil, model.NewValidationError("account_id", "账户不存在")
	}
	holdings := s.store.ListHoldings()
	var result []HoldingSummary
	for _, h := range holdings {
		if h.AccountID != accountID {
			continue
		}
		fund, err := s.store.GetFund(h.FundID)
		if err != nil {
			continue
		}
		result = append(result, HoldingSummary{
			FundID:      h.FundID,
			FundName:    fund.Name,
			FundCode:    fund.Code,
			TotalShares: h.TotalShares,
			Cost:        h.Cost,
			MarketValue: h.MarketValue,
			Profit:      h.MarketValue - h.Cost,
		})
	}
	return result, nil
}

func (s *Service) GetAccountProfit(accountID string) (*AccountProfit, error) {
	account, err := s.store.GetAccount(accountID)
	if err != nil {
		return nil, model.NewValidationError("account_id", "账户不存在")
	}
	holdings := s.store.ListHoldings()
	var totalCost, totalMarket int64
	for _, h := range holdings {
		if h.AccountID != accountID {
			continue
		}
		totalCost += h.Cost
		totalMarket += h.MarketValue
	}
	return &AccountProfit{
		AccountID:   accountID,
		Owner:       account.Owner,
		TotalCost:   totalCost,
		TotalMarket: totalMarket,
		TotalProfit: totalMarket - totalCost,
	}, nil
}

func (s *Service) GetFundTypeStats() []FundTypeStat {
	funds := s.store.ListFunds()
	m := make(map[string]*FundTypeStat)
	for _, f := range funds {
		if _, ok := m[f.Type]; !ok {
			m[f.Type] = &FundTypeStat{Type: f.Type}
		}
		m[f.Type].Count++
		m[f.Type].TotalNav += f.Nav
	}
	result := make([]FundTypeStat, 0, len(m))
	for _, v := range m {
		result = append(result, *v)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Type < result[j].Type
	})
	return result
}

func (s *Service) GetFundRiskStats() []FundRiskStat {
	funds := s.store.ListFunds()
	m := make(map[int]*FundRiskStat)
	for _, f := range funds {
		if _, ok := m[f.RiskLevel]; !ok {
			m[f.RiskLevel] = &FundRiskStat{RiskLevel: f.RiskLevel}
		}
		m[f.RiskLevel].Count++
	}
	result := make([]FundRiskStat, 0, len(m))
	for _, v := range m {
		result = append(result, *v)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].RiskLevel < result[j].RiskLevel
	})
	return result
}

func (s *Service) GetTopFunds(n int) []TopFund {
	if n <= 0 {
		n = 10
	}
	txns := s.store.ListTransactions()
	m := make(map[string]int64)
	for _, t := range txns {
		if t.Type == model.TransactionTypeInvest && t.Status == model.TransactionStatusConfirmed {
			m[t.FundID] += t.Amount
		}
	}
	type pair struct {
		fundID string
		amount int64
	}
	var pairs []pair
	for k, v := range m {
		pairs = append(pairs, pair{fundID: k, amount: v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].amount > pairs[j].amount
	})
	var result []TopFund
	for i, p := range pairs {
		if i >= n {
			break
		}
		fund, err := s.store.GetFund(p.fundID)
		if err != nil {
			continue
		}
		result = append(result, TopFund{
			FundID:        p.fundID,
			FundName:      fund.Name,
			FundCode:      fund.Code,
			TotalInvested: p.amount,
		})
	}
	return result
}
