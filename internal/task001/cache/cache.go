package cache

import (
	"sync"

	"fundinvest/internal/task001/builder"
)

type Cache struct {
	mu    sync.RWMutex
	items map[string]*builder.Portfolio
}

func New() *Cache { return &Cache{items: make(map[string]*builder.Portfolio)} }

func (c *Cache) Put(p *builder.Portfolio) error {
	if p == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[p.ID] = p
	return nil
}

func (c *Cache) Get(id string) (*builder.Portfolio, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	p, ok := c.items[id]
	return p, ok
}
