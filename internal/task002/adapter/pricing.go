package adapter

import (
	"context"
	"sync/atomic"
)

type Pricing struct{ calls atomic.Int64 }

func (p *Pricing) Fetch(ctx context.Context, fund string) (int64, error) {
	p.calls.Add(1)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return 12345, nil
	}
}

func (p *Pricing) Calls() int64 { return p.calls.Load() }
