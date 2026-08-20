package flow

import (
	"errors"
	"fundinvest/internal/task010/audit"
	"fundinvest/internal/task010/resource"
	"fundinvest/internal/task010/service"
	"fundinvest/internal/task010/transaction"
	"testing"
)

func TestBatchReleasesEachHandleAndPreservesBusinessFailure(t *testing.T) {
	pool, registry, log := &resource.Pool{Limit: 2}, &transaction.Registry{}, &audit.Log{}
	batch := &service.Batch{Pool: pool, Registry: registry, Audit: log}
	err := batch.Process([]string{"fund-a", "fund-b", "fund-c", "blocked-fund"})
	if !errors.Is(err, transaction.ErrInvalidFund) {
		t.Errorf("business failure was replaced: %v", err)
	}
	if pool.Open != 0 || pool.Peak != 1 {
		t.Errorf("statement handles accumulated: open=%d peak=%d", pool.Open, pool.Peak)
	}
	if got := registry.Committed(); len(got) != 3 {
		t.Errorf("batch stopped before valid items committed: %#v", got)
	}
	if entries := log.Entries(); len(entries) != 3 || entries[0] != "fund-a" || entries[1] != "fund-b" || entries[2] != "fund-c" {
		t.Errorf("audit and committed items diverged: %#v", entries)
	}
}
