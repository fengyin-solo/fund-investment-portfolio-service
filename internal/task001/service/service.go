package service

import (
	"fundinvest/internal/task001/builder"
	"fundinvest/internal/task001/cache"
)

type Manager struct {
	builder builder.Builder
	cache   *cache.Cache
}

func New(c *cache.Cache) *Manager { return &Manager{cache: c} }

func (m *Manager) Prepare(id string, interrupt bool) error {
	portfolio, err := m.builder.Build(id, interrupt)
	if portfolio != nil {
		_ = m.cache.Put(portfolio)
	}
	return err
}

func (m *Manager) Lookup(id string) (*builder.Portfolio, bool) {
	return m.cache.Get(id)
}
