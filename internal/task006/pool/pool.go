package pool

import (
	"sync"

	"fundinvest/internal/task006/identity"
)

type Pool struct {
	mu     sync.Mutex
	cached *identity.Lease
}

func (p *Pool) Acquire() *identity.Lease {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cached != nil {
		value := p.cached
		p.cached = nil
		return value
	}
	return &identity.Lease{}
}

func (p *Pool) Release(value *identity.Lease) {
	p.mu.Lock()
	p.cached = value
	p.mu.Unlock()
}
