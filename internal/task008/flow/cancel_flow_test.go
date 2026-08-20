package flow

import (
	"context"
	"testing"
	"time"

	"fundinvest/internal/task008/client"
	"fundinvest/internal/task008/service"
	"fundinvest/internal/task008/tracker"
)

func TestCancelStopsQuoteRetryAndAllowsShutdown(t *testing.T) {
	permits := make(chan struct{}, 3)
	started := make(chan int, 3)
	group := &tracker.Group{}
	downstream := &client.Client{Permit: permits, Started: started}
	quotes := &service.Quotes{Client: downstream, Tasks: group}
	ctx, cancel := context.WithCancel(context.Background())
	quotes.Start(ctx)
	if call := <-started; call != 1 {
		t.Fatalf("unexpected first call: %d", call)
	}
	permits <- struct{}{}
	if call := <-started; call != 2 {
		t.Fatalf("unexpected second call: %d", call)
	}
	cancel()
	waitCtx, stopWait := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer stopWait()
	if err := group.Wait(waitCtx); err != nil {
		t.Errorf("shutdown waited for canceled quote work: %v", err)
	}
	permits <- struct{}{}
	select {
	case call := <-started:
		t.Errorf("quote retry %d started after cancellation", call)
	case <-time.After(5 * time.Millisecond):
	}
	if calls := downstream.Calls(); calls != 2 {
		t.Errorf("quote retries continued after cancellation: %d calls", calls)
	}
	permits <- struct{}{}
}
