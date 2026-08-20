package scope

import (
	"context"
	"sync"
)

type Router struct {
	mu    sync.Mutex
	first context.Context
}

func New() *Router { return &Router{} }

func (r *Router) ForRequest(ctx context.Context) context.Context {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.first == nil {
		r.first = ctx
	}
	return r.first
}
