package flow

import (
	"fundinvest/internal/task007/cache"
	"fundinvest/internal/task007/effects"
	"fundinvest/internal/task007/model"
	"fundinvest/internal/task007/service"
	"fundinvest/internal/task007/worker"
	"testing"
)

func TestRetryKeepsSingleEffectAndRejectsDelayedState(t *testing.T) {
	store, client := cache.New(), effects.New()
	w := worker.Worker{Cache: store}
	original := model.New("plan-7")
	completed, err := service.Retry(original, client, w)
	if err != nil {
		t.Fatalf("retry did not recover: %v", err)
	}
	w.Delayed(original)
	if client.SideEffects != 1 {
		t.Errorf("settlement side effect ran %d times", client.SideEffects)
	}
	cached := store.Get(original.ID)
	if completed.State != "succeeded" || cached.State != "succeeded" || cached.Version != completed.Version {
		t.Errorf("delayed callback rolled final state back: completed=%#v cached=%#v", completed, cached)
	}
}
