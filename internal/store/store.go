// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"fundinvest/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	CreateFund(f *model.Fund) error
	GetFund(id string) (*model.Fund, error)
	GetFundByCode(code string) (*model.Fund, error)
	ListFunds() []*model.Fund
	UpdateFund(f *model.Fund) error
	DeleteFund(id string) error

	CreateAccount(a *model.Account) error
	GetAccount(id string) (*model.Account, error)
	ListAccounts() []*model.Account
	UpdateAccount(a *model.Account) error
	DeleteAccount(id string) error

	CreatePlan(p *model.Plan) error
	GetPlan(id string) (*model.Plan, error)
	ListPlans() []*model.Plan
	UpdatePlan(p *model.Plan) error
	DeletePlan(id string) error

	CreateTransaction(t *model.Transaction) error
	GetTransaction(id string) (*model.Transaction, error)
	ListTransactions() []*model.Transaction
	UpdateTransaction(t *model.Transaction) error
	DeleteTransaction(id string) error

	CreateHolding(h *model.Holding) error
	GetHolding(id string) (*model.Holding, error)
	GetHoldingByAccountAndFund(accountID, fundID string) (*model.Holding, error)
	ListHoldings() []*model.Holding
	UpdateHolding(h *model.Holding) error
	DeleteHolding(id string) error

	CreateNavHistory(n *model.NavHistory) error
	GetNavHistory(id string) (*model.NavHistory, error)
	ListNavHistory() []*model.NavHistory
	UpdateNavHistory(n *model.NavHistory) error
	DeleteNavHistory(id string) error

	CreateDividend(d *model.Dividend) error
	GetDividend(id string) (*model.Dividend, error)
	ListDividends() []*model.Dividend
	UpdateDividend(d *model.Dividend) error
	DeleteDividend(id string) error
}
