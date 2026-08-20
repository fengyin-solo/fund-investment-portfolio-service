package tracker

import (
	"context"
	"sync"
)

type Group struct{ wg sync.WaitGroup }

func (g *Group) Go(run func() error) {
	g.wg.Add(1)
	go func() {
		if run() == nil {
			g.wg.Done()
		}
	}()
}

func (g *Group) Wait(ctx context.Context) error {
	done := make(chan struct{})
	go func() { g.wg.Wait(); close(done) }()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
