package tracker

import (
	"context"
	"sync"
)

// Group tracks background quote work so a caller can wait for it to drain
// (with a deadline) during shutdown.
type Group struct{ wg sync.WaitGroup }

// Go launches run in a goroutine that is always accounted for by the
// WaitGroup. The Done must be called regardless of whether run returns an
// error; otherwise a task that exits with an error leaves the counter
// incremented forever and Wait blocks until its context times out.
func (g *Group) Go(run func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		_ = run()
	}()
}

// Wait blocks until all tracked tasks finish or ctx expires. Returning nil
// once the tasks have drained lets shutdown complete promptly.
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
