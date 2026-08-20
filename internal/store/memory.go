package store

import (
	"sync"

	"fundinvest/internal/model"
)

type MemoryStore struct {
	mu            sync.RWMutex
	funds         map[string]*model.Fund
	accounts      map[string]*model.Account
	plans         map[string]*model.Plan
	transactions  map[string]*model.Transaction
	holdings      map[string]*model.Holding
	navHistory    map[string]*model.NavHistory
	dividends     map[string]*model.Dividend
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		funds:        make(map[string]*model.Fund),
		accounts:     make(map[string]*model.Account),
		plans:        make(map[string]*model.Plan),
		transactions: make(map[string]*model.Transaction),
		holdings:     make(map[string]*model.Holding),
		navHistory:   make(map[string]*model.NavHistory),
		dividends:    make(map[string]*model.Dividend),
	}
}

var _ Store = (*MemoryStore)(nil)
